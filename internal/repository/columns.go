package repository

//	var orderColumns = []string{
//		"order_uid",
//		"track_number",
//		"entry",
//		"locale",
//		"internal_signature",
//		"customer_id",
//		"delivery_service",
//		"shardkey",
//		"sm_id",
//		"date_created",
//		"oof_shard",
//	}
var orderColumns = "order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard"

// var deliveryColumns = []string{
// 	"name",
// 	"phone",
// 	"zip",
// 	"city",
// 	"address",
// 	"region",
// 	"email",
// }

var deliveryColumns = "name, phone, zip, city, address, region, email"

var deliveryColumnsWithOrderUID = "order_uid, name, phone, zip, city, address, region, email"

// var paymentColumns = []string{
// 	"transaction",
// 	"request_id",
// 	"currency",
// 	"provider",
// 	"amount",
// 	"payment_dt",
// 	"bank",
// 	"delivery_cost",
// 	"goods_total",
// 	"custom_fee",
// }

var paymentColumns = "transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee"

var paymentColumnsWithOrderUID = "order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee"

//	var itemColumns = []string{
//		"chrt_id",
//		"track_number",
//		"price",
//		"rid",
//		"name",
//		"sale",
//		"size",
//		"total_price",
//		"nm_id",
//		"brand",
//		"status",
//	}
var itemColumns = "chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status"

var itemColumnsWithOrderUID = "order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status"
