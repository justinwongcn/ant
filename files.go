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

// FileUploader 文件上传处理器
// 提供安全的文件上传功能，支持自定义上传目标路径
type FileUploader struct {
	// FileField 表单中文件字段的名称
	FileField string
	// DstPathFunc 根据上传的文件信息确定目标存储路径的函数
	DstPathFunc func(fh *multipart.FileHeader) string
	// FileNameFunc 生成文件名的函数，如果为nil则使用原始文件名
	FileNameFunc func(originalName string) string
}

// Handle 实现文件上传处理逻辑
// 返回值: 返回处理上传请求的HandleFunc
// 注意：
// 1. 自动创建目标目录
// 2. 返回上传结果和文件大小信息
// 3. 处理各类错误场景并返回适当的HTTP状态码
// 4. 支持自定义文件名生成策略，避免文件重名
func (f *FileUploader) Handle() HandleFunc {
	return func(ctx *Context) {
		src, fileHeader, err := ctx.Req.FormFile(f.FileField)
		if err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("上传失败，未找到文件")
			return
		}
		defer src.Close()

		// 生成文件名
		originalName := fileHeader.Filename
		fileName := originalName
		if f.FileNameFunc != nil {
			fileName = f.FileNameFunc(originalName)
		}

		// 使用新的文件名创建FileHeader
		newFileHeader := &multipart.FileHeader{
			Filename: fileName,
			Size:     fileHeader.Size,
			Header:   fileHeader.Header,
		}

		// 确保目标目录存在
		dstPath := f.DstPathFunc(newFileHeader)
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

// FileDownloader 文件下载处理器
// 提供安全的文件下载功能，支持防止目录遍历攻击
type FileDownloader struct {
	// Dir 文件下载的根目录
	Dir string
}

// Handle 实现文件下载处理逻辑
// 返回值: 返回处理下载请求的HandleFunc
// 注意：
// 1. 自动处理文件不存在、权限错误等异常情况
// 2. 防止目录遍历和路径穿越攻击
// 3. 设置正确的Content-Type和Content-Disposition头
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

// StaticResourceHandler 静态资源处理器
// 提供高性能的静态资源服务，支持文件缓存和自定义Content-Type
type StaticResourceHandler struct {
	// dir 静态资源的根目录
	dir string
	// pathPrefix 静态资源的URL路径前缀
	pathPrefix string
	// extensionContentTypeMap 文件扩展名到Content-Type的映射
	extensionContentTypeMap map[string]string
	// cache 文件内容缓存
	cache *lru.Cache
	// maxFileSize 可缓存的最大文件大小
	maxFileSize int
}

// fileCacheItem 文件缓存项
// 用于在内存中缓存静态资源文件的内容和元数据
type fileCacheItem struct {
	// fileName 文件名
	fileName string
	// fileSize 文件大小（字节）
	fileSize int
	// contentType 文件的Content-Type
	contentType string
	// data 文件内容
	data []byte
	// modTime 文件修改时间戳
	modTime int64
}

// StaticResourceHandlerOption 定义静态资源处理器的配置选项函数类型
type StaticResourceHandlerOption func(*StaticResourceHandler)

// NewStaticResourceHandler 创建新的静态资源处理器
// dir: 静态资源的根目录
// pathPrefix: 静态资源的URL路径前缀
// options: 可选的配置选项
// 返回值: 配置完成的StaticResourceHandler实例
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

// Handle 处理静态资源请求
// ctx: 请求上下文
// 注意：
// 1. 支持从缓存中快速返回资源
// 2. 自动设置适当的Content-Type
// 3. 处理各类错误场景
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

// readFileFromData 从缓存中读取文件数据
// fileName: 要读取的文件名
// 返回值:
// - *fileCacheItem: 缓存的文件项
// - bool: 是否找到缓存项
func (h *StaticResourceHandler) readFileFromData(fileName string) (*fileCacheItem, bool) {
	if h.cache == nil {
		return nil, false
	}
	if item, ok := h.cache.Get(fileName); ok {
		return item.(*fileCacheItem), true
	}
	return nil, false
}

// writeItemAsResponse 将缓存项写入HTTP响应
// item: 要写入的缓存项
// writer: HTTP响应写入器
// 注意：设置适当的HTTP头部，包括缓存控制
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

// cacheFile 将文件缓存到内存中
// item: 要缓存的文件项
// 注意：只有文件大小小于maxFileSize时才会被缓存
func (h *StaticResourceHandler) cacheFile(item *fileCacheItem) {
	if h.cache != nil && item.fileSize < h.maxFileSize {
		h.cache.Add(item.fileName, item)
	}
}

// WithFileCache 创建启用文件缓存的配置选项
// maxFileSizeThreshold: 可缓存的最大文件大小（字节）
// maxCacheFileCnt: 缓存中可存储的最大文件数量
// 返回值: StaticResourceHandlerOption配置函数
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

// WithMoreExtension 创建扩展Content-Type映射的配置选项
// extMap: 要添加的扩展名到Content-Type的映射
// 返回值: StaticResourceHandlerOption配置函数
func WithMoreExtension(extMap map[string]string) StaticResourceHandlerOption {
	return func(h *StaticResourceHandler) {
		maps.Copy(h.extensionContentTypeMap, extMap)
	}
}

// getFileExt 获取文件名中的扩展名
// name: 完整的文件名
// 返回值: 文件扩展名（不包含点号），如果没有扩展名则返回空字符串
func getFileExt(name string) string {
	index := strings.LastIndex(name, ".")
	if index == len(name)-1 {
		return ""
	}
	return name[index+1:]
}
