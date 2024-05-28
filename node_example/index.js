const express = require("express");
const crypto = require("crypto");
var nodeBase64 = require('nodejs-base64-converter');
var request = require("request");








const app = express();



app.use(express.json());
app.use(express.urlencoded({ extended: true }));


const  handlePayment = async(req,res) => {

	const {data} = req.body;

	console.log(data)

	var test_basket = [];


	data.map((item) => {
		test_basket.push([
			item.name,
			item.price,
			item.amount
		])
	})


	 var merchant_id = "272972";
    var merchant_key = "fKbPJLPc2UtyUNet";

    var merchant_salt = "gTp58ZLH7Wjhxrhi";

    var user_basket = nodeBase64.encode(test_basket);

    var merchant_oid = "test123";

    var max_installment = 0;
    var no_installment = '0'  // Taksit yapılmasını istemiyorsanız, sadece tek çekim sunacaksanız 1 yapın.
    var user_ip ="78.163.140.60";
    var email = "test@test.com";// Müşterinizin sitenizde kayıtlı veya form vasıtasıyla aldığınız eposta adresi.
    var payment_amount = 100*10; // Tahsil edilecek tutar. 9.99 için 9.99 * 100 = 999 gönderilmelidir.
    var currency = 'TL';
    var test_mode = '0'; // Mağaza canlı modda iken test işlem yapmak için 1 olarak gönderilebilir.
    var user_name = "test"; // Müşterinizin sitenizde kayıtlı veya form aracılığıyla aldığınız ad ve soyad bilgisi
    var user_address = "test_address"; // Müşterinizin sitenizde kayıtlı veya form aracılığıyla aldığınız adres bilgisi
    var user_phone = "123235235432"; // Müşterinizin sitenizde kayıtlı veya form aracılığıyla aldığınız telefon bilgisi

    // Başarılı ödeme sonrası müşterinizin yönlendirileceği sayfa
    // Bu sayfa siparişi onaylayacağınız sayfa değildir! Yalnızca müşterinizi bilgilendireceğiniz sayfadır!
    var merchant_ok_url = 'https://apps.uniqgene.com/checkout';
    // Ödeme sürecinde beklenmedik bir hata oluşması durumunda müşterinizin yönlendirileceği sayfa
    // Bu sayfa siparişi iptal edeceğiniz sayfa değildir! Yalnızca müşterinizi bilgilendireceğiniz sayfadır!
    var merchant_fail_url = 'https://apps.uniqgene.com/';
    var timeout_limit = 30; // İşlem zaman aşımı süresi - dakika cinsinden
    var debug_on = 1; // Hata mesajlarının ekrana basılması için entegrasyon ve test sürecinde 1 olarak bırakın. Daha sonra 0 yapabilirsiniz.
    var lang = 'tr'; // Türkçe için tr veya İngilizce için en gönderilebilir. Boş gönderilirse tr geçerli olur.

    var hashSTR = `${merchant_id}${user_ip}${merchant_oid}${email}${payment_amount}${user_basket}${no_installment}${max_installment}${currency}${test_mode}`;

    var paytr_token = hashSTR + merchant_salt;

    var token = crypto.createHmac('sha256', merchant_key).update(paytr_token).digest('base64');

	console.log(token)

    var options = {
        method: 'POST',
        url: 'https://www.paytr.com/odeme/api/get-token',
        headers:
            { 'content-type': 'application/x-www-form-urlencoded' },
        formData: {
            merchant_id: merchant_id,
            merchant_key: merchant_key,
            merchant_salt: merchant_salt,
            email: email,
            payment_amount: payment_amount,
            merchant_oid: merchant_oid,
            user_name: user_name,
            user_address: user_address,
            user_phone: user_phone,
            merchant_ok_url: merchant_ok_url,
            merchant_fail_url: merchant_fail_url,
            user_basket: user_basket,
            user_ip: user_ip,
            timeout_limit: timeout_limit,
            debug_on: debug_on,
            test_mode: test_mode,
            lang: lang,
            no_installment: no_installment,
            max_installment: max_installment,
            currency: currency,
            paytr_token: token,


        }
    };

    request(options, function (error, response, body) {
        if (error) throw new Error(error);
        var res_data = JSON.parse(body);

        if (res_data.status == 'success') {
            return res.send( { iframetoken: res_data.token });
        } else {

           return res.send(body);
        }


    });





}

app.post("/test-payment", handlePayment)


app.listen(3001, () => {
    console.log("Server is running on port 3001");
});








