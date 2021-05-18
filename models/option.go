//Package key redis的key包装
package models

//Option 设置key行为的选项
//@attribute MaxTTL time.Duration 为0则不设置过期
//@attribute AutoRefresh string 需要为crontab格式的字符串,否则不会自动定时刷新
type Options struct {
	ShowSQL      bool
	MaxIdle      int
	MaxOpenConns int
}

// Option configures how we set up the connection.
type Option interface {
	Apply(*Options)
}

// func (emptyOption) apply(*Options) {}
type funcOption struct {
	f func(*Options)
}

func (fo *funcOption) Apply(do *Options) {
	fo.f(do)
}

func newFuncOption(f func(*Options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

//WithMaxTTL 设置最大过期时间
func ShowSQL() Option {
	return newFuncOption(func(o *Options) {
		o.ShowSQL = true
	})
}

//WithMaxIdle 设置最长空置时间
func WithMaxIdle(MaxIdle int) Option {
	return newFuncOption(func(o *Options) {
		o.MaxIdle = MaxIdle
	})
}

//WithMaxOpenConns 设置最大打开连接数
func WithMaxOpenConns(MaxOpenConns int) Option {
	return newFuncOption(func(o *Options) {
		o.MaxOpenConns = MaxOpenConns
	})
}
