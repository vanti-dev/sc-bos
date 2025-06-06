package api

var statusCodes = map[string]string{
	"0x02401000": "No AppKey is configured. Enter the correct AppKey in the\nrequest.",
	"0x02401001": "The partner of AppKey does not exist. Check the AppKey in the	request.",
	"0x02401002": "No signature is configured. Enter the correct signature in the request.",
	"0x02401003": "Invalid signature. Check the signature in the request.",
	"0x02401004": "Token authentication failed. Check the token.",
	"0x02401005": "No token is configured. Enter the token.",
	"0x02401006": "Token exception. Check the token.",
	"0x02401007": "No permission. Please contact the administrator to apply for permissions.",
	"0x02401008": "Authentication exception. Check the gateway service.",
	"0x02401009": "Maximum API calling attempts reached. Please contact the administrator to apply for adding access attempts.",
	"0x0240100a": "Parameter conversion exception. Check the API parameters.",
	"0x0240100b": "Calling statistics exception. Check the gateway.",

	"0x00072001": "The required parameters are not configured. Set the required parameters in the request.",
	"0x00072002": "Invalid parameter value range.",
	"0x00072003": "Invalid parameter value format.",
	"0x00072004": "The response message is too long. Set the page size in the request.",

	"0x00052101": "Highest service performance reached. Try again later.",
	"0x00052102": "Service error. Try again later.",
	"0x00052103": "Service response timed out. Try again later.",
	"0x00052104": "Service is not available. Try again after restoring the service",

	"0x00072201": "No permission for resource access. Please contact the	administrator to apply for permissions.",
	"0x00072202": "The resource does not exist. Enter the correct resource No. in the request.",
	"0x00072203": "Maximum number of Licenses reached. Check the License information from the administrator.",
	"0x00072204": "No permission for this function. Check the License information from the administrator.",

	"0x00052301": "Unknown error.",
}
