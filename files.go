package ant

import (
	"fmt"
	"io"
	"log"
	"maps"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

type FileUploader struct {
	FileField   string
	DstPathFunc func(fh *multipart.FileHeader) string
}

func (f *FileUploader) Handle() HandleFunc {
	return func(ctx *Context) {
		src, fileHeader, err := ctx.Req.FormFile(f.FileField)
		if err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("上传失败，未找到文件")
			return
		}
		defer src.Close()

		// 确保目标目录存在
		dstPath := f.DstPathFunc(fileHeader)
		if err = os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("创建目录失败")
			ctx.Resp.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("创建文件失败")
			ctx.Resp.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer dst.Close()

		written, err := io.Copy(dst, src)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("保存文件失败")
			ctx.Resp.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = fmt.Appendf(nil, "上传成功，文件大小: %d bytes", written)
	}
}

type FileDownloader struct {
	Dir string
}

func (f *FileDownloader) Handle() HandleFunc {
	return func(ctx *Context) {
		fileName, err := ctx.QueryValue("file").String()
		if err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("未指定文件名")
			ctx.Resp.WriteHeader(http.StatusBadRequest)
			return
		}

		// 清理文件路径，防止目录遍历攻击
		cleanPath := filepath.Clean(fileName)
		if strings.Contains(cleanPath, "..") {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("非法的文件路径")
			ctx.Resp.WriteHeader(http.StatusBadRequest)
			return
		}

		// 使用filepath.Base确保路径限制在目标目录内，防止绝对路径攻击
		filePath := filepath.Join(f.Dir, filepath.Base(cleanPath))
		info, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				ctx.RespStatusCode = http.StatusNotFound
				ctx.RespData = []byte("文件不存在")
				ctx.Resp.WriteHeader(http.StatusNotFound)
			} else if os.IsPermission(err) {
				// 添加对权限错误的特殊处理，使用403 Forbidden状态码
				ctx.RespStatusCode = http.StatusForbidden
				ctx.RespData = []byte("没有访问权限")
				ctx.Resp.WriteHeader(http.StatusForbidden)
			} else {
				ctx.RespStatusCode = http.StatusInternalServerError
				ctx.RespData = []byte("读取文件失败")
				ctx.Resp.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		// 检查是否是目录
		if info.IsDir() {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("打开文件失败")
			ctx.Resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		file, err := os.Open(filePath)
		if err != nil {
			// 区分权限错误和其他错误
			if os.IsPermission(err) {
				ctx.RespStatusCode = http.StatusForbidden
				ctx.RespData = []byte("读取文件失败")
				ctx.Resp.WriteHeader(http.StatusForbidden)
			} else {
				ctx.RespStatusCode = http.StatusInternalServerError
				ctx.RespData = []byte("打开文件失败")
				ctx.Resp.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		defer file.Close()

		// 设置响应头
		header := ctx.Resp.Header()
		header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(cleanPath)))
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Length", fmt.Sprintf("%d", info.Size()))

		// 设置响应状态码
		ctx.RespStatusCode = http.StatusOK
		ctx.Resp.WriteHeader(http.StatusOK)
		_, err = io.Copy(ctx.Resp, file)
		if err != nil {
			log.Printf("发送文件失败: %v", err)
		}
	}
}

type StaticResourceHandler struct {
	dir                     string
	pathPrefix              string
	extensionContentTypeMap map[string]string
	cache                   *lru.Cache
	maxFileSize             int
}

type fileCacheItem struct {
	fileName    string
	fileSize    int
	contentType string
	data        []byte
	modTime     int64
}

type StaticResourceHandlerOption func(*StaticResourceHandler)

func NewStaticResourceHandler(dir, pathPrefix string, options ...StaticResourceHandlerOption) *StaticResourceHandler {
	h := &StaticResourceHandler{
		dir:        dir,
		pathPrefix: pathPrefix,
		extensionContentTypeMap: map[string]string{
			"html": "text/html; charset=utf-8",
			"css":  "text/css; charset=utf-8",
			"js":   "application/javascript",
			"json": "application/json",
			"jpeg": "image/jpeg",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"gif":  "image/gif",
			"ico":  "image/x-icon",
			"svg":  "image/svg+xml",
			"pdf":  "application/pdf",
			"txt":  "text/plain",
		},
	}

	for _, opt := range options {
		opt(h)
	}
	return h
}

func (h *StaticResourceHandler) Handle(ctx *Context) {
	// 获取请求路径中的文件名
	req, err := ctx.PathValue("file").String()
	if err != nil {
		ctx.RespStatusCode = http.StatusBadRequest
		ctx.RespData = []byte("未指定文件名")
		return
	}

	// 检查是否包含非法路径
	if strings.Contains(req, "..") {
		ctx.RespStatusCode = http.StatusBadRequest
		ctx.RespData = []byte("未指定文件名")
		return
	}

	// 从数据中读取文件内容
	item, ok := h.readFileFromData(req)
	if ok {
		// 如果文件存在，则从缓存中写入响应并返回
		log.Printf("从缓存中读取数据...")
		h.writeItemAsResponse(item, ctx.Resp)
		return
	}

	// 拼接文件路径
	path := filepath.Join(h.dir, req)
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		// 如果文件打开失败，则返回内部服务器错误状态码
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte("打开文件失败")
		return
	}
	defer file.Close()

	// 获取文件扩展名
	ext := getFileExt(file.Name())
	// 根据扩展名获取对应的 content type
	t, ok := h.extensionContentTypeMap[ext]
	if !ok {
		// 如果扩展名对应的 content type 不存在，则返回Bad Request状态码
		ctx.RespStatusCode = http.StatusBadRequest
		ctx.RespData = []byte("不支持的文件类型")
		return
	}

	// 读取文件内容
	data, err := io.ReadAll(file)
	if err != nil {
		// 如果读取文件失败，则返回内部服务器错误状态码
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte("读取文件失败")
		return
	}

	// 创建 fileCacheItem 对象并设置属性值
	item = &fileCacheItem{
		fileName:    req,
		fileSize:    len(data),
		contentType: t,
		data:        data,
		modTime:     time.Now().Unix(),
	}

	// 将文件缓存到内存中
	h.cacheFile(item)
	// 将 fileCacheItem 对象写入响应并返回
	ctx.RespStatusCode = http.StatusOK
	h.writeItemAsResponse(item, ctx.Resp)
}

func (h *StaticResourceHandler) readFileFromData(fileName string) (*fileCacheItem, bool) {
	if h.cache == nil {
		return nil, false
	}
	if item, ok := h.cache.Get(fileName); ok {
		return item.(*fileCacheItem), true
	}
	return nil, false
}

func (h *StaticResourceHandler) writeItemAsResponse(item *fileCacheItem, writer http.ResponseWriter) {
	header := writer.Header()
	header.Set("Content-Type", item.contentType)
	header.Set("Content-Length", fmt.Sprintf("%d", item.fileSize))
	header.Set("Last-Modified", fmt.Sprintf("%d", item.modTime))
	header.Set("Cache-Control", "public, max-age=31536000")
	writer.WriteHeader(http.StatusOK)
	_, err := writer.Write(item.data)
	if err != nil {
		log.Printf("写入响应失败: %v", err)
	}
}

func (h *StaticResourceHandler) cacheFile(item *fileCacheItem) {
	if h.cache != nil && item.fileSize < h.maxFileSize {
		h.cache.Add(item.fileName, item)
	}
}

func WithFileCache(maxFileSizeThreshold int, maxCacheFileCnt int) StaticResourceHandlerOption {
	return func(h *StaticResourceHandler) {
		cache, err := lru.New(maxCacheFileCnt)
		if err != nil {
			log.Printf("创建缓存失败: %v", err)
			return
		}
		h.maxFileSize = maxFileSizeThreshold
		h.cache = cache
	}
}

func WithMoreExtension(extMap map[string]string) StaticResourceHandlerOption {
	return func(h *StaticResourceHandler) {
		maps.Copy(h.extensionContentTypeMap, extMap)
	}
}

// getFileExt 函数获取文件名中的文件扩展名
func getFileExt(name string) string {
	index := strings.LastIndex(name, ".")
	if index == len(name)-1 {
		return ""
	}
	return name[index+1:]
}
