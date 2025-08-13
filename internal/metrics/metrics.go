package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	LoginSuccessTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "login_success_total",
		Help: "Количество успешных логинов",
	})

	LoginFailedTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "login_failed_total",
		Help: "Количество неудачных логинов",
	})

	RegisterTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "register_total",
		Help: "Количество регистраций пользователей",
	})

	ImagesUploaded = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "images_uploaded_total",
		Help: "Количество загруженных изображений",
	})

	ProductsCreated = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "products_created_total",
		Help: "Количество созданных товаров",
	})

	ProductsFetched = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "products_fetched_total",
		Help: "Количество запросов на получение списка товаров",
	})

	ProductFetched = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "product_fetched_total",
		Help: "Количество запросов на получение товара по slug",
	})
	ProductsUpdated = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "products_updated_total",
		Help: "Количество обновлённых товаров",
	})

	ProductsDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "products_deleted_total",
		Help: "Количество удалённых товаров",
	})

	OAuthLoginTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "oauth_login_total",
		Help: "Количество логинов через OAuth-провайдеров",
	})
)

func Init() {
	prometheus.MustRegister(
		LoginSuccessTotal,
		LoginFailedTotal,
		RegisterTotal,
		ImagesUploaded,
		ProductsCreated,
		ProductsFetched,
		ProductFetched,
		ProductsUpdated,
		ProductsDeleted,
		OAuthLoginTotal,
	)
}
