
// const client = new EventServiceClient('http://13.232.13.253:8080');
const {Events} = require("./proto/eventservice_pb");
const {Header} = require("./proto/eventservice_pb");
const {EventServiceClient} = require("./proto/eventservice_grpc_web_pb");
const client = new EventServiceClient('http://localhost:8080');

const header = new Header();
header.setOs("PWA");
const req = new Events().setHeader(header);

client.dispatch(req, {}, (err, response) => {
    if (err) {
        console.log(`Unexpected error: code = ${err.code}` +
            `, message = "${err.message}"`);
    } else {
        console.log(response.getStatus());
    }
});