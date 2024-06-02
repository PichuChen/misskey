
use actix_web::{web, App, HttpServer, Responder, HttpResponse, HttpRequest};
use reqwest;

const MISSKEY_NODE_BACKEND: &str = "http://localhost:3000";
const BACKEND2_LISTEN_PORT: u16 = 3002;

async fn index(r: HttpRequest, bytes: web::Bytes) -> impl Responder{
	let request_url = r.path();
	let request_method = r.method();
	let request_headers = r.head().headers();

	println!("Request received: method: {}, url: {}", request_method, request_url);

	// send request to misskey-node
	let method_str = request_method.to_string();
	let method = reqwest::Method::from_bytes(method_str.as_bytes()).unwrap();

	let mut headers = reqwest::header::HeaderMap::new();

	for (key, value) in request_headers.iter() {
		let header_key: reqwest::header::HeaderName = reqwest::header::HeaderName::from_bytes(key.to_string().as_bytes()).unwrap();
		let header_value =  reqwest::header::HeaderValue::from_str(value.to_str().unwrap()).unwrap();

		println!("Request headers {}:{}", header_key, header_value.to_str().unwrap());
		headers.append(header_key, header_value);
	}
	let body = String::from_utf8(bytes.to_vec()).unwrap();
	let misskey_url = format!("{}{}", MISSKEY_NODE_BACKEND, request_url);
	let response = call_node_backend(method, misskey_url, headers, body).await;

	match response {
		Ok((res, h)) => {
			// println!("response: {}", res);
			let mut response = HttpResponse::Ok();
			for (key, value) in h.iter() {
				let header_key = key.as_str();
				let header_value = value.to_str().unwrap();
				// println!("Response headers {}: {}", header_key, header_value);
			  response.insert_header((header_key, header_value));
			}
			return response
				.body(res);
		},
		Err(e) => {
			println!("error: {}", e);
			return HttpResponse::InternalServerError()
				.body("Internal Server Error");
		}
	}
}


async fn call_node_backend(method :reqwest::Method,
		url: String,
		headers: reqwest::header::HeaderMap,
		body: String
	) -> Result<(String, reqwest::header::HeaderMap), reqwest::Error>{
	let m = method.clone();
	let url2 = url.clone();
	println!("method: {}, url: {}", method, url);

	let res = reqwest::Client::new()
		.request(m, url2)
		.headers(headers)
		.body(body)
		.send()
		.await?;
	// let res = reqwest::get("http://httpbin.org/get").await?;
	println!("Status: {}", res.status());
	println!("Response headers:\n{:#?}", res.headers());

	let response_headers = res.headers().clone();
	let reponse_body = res.text().await?;
	// println!("Body:\n{}", reponseBody);

	Ok((reponse_body, response_headers))
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
	// Create a http server
	let server = HttpServer::new(|| {
		App::new()
			// Handle all request
			.default_service(web::to(index))
			// TODO: multi-thread?
	});

	// Bind the server to the local address
	server.bind(("127.0.0.1", BACKEND2_LISTEN_PORT))?.run().await
}


