package methods

const SignPing = "client.ping"
const SignPingRET = "success"

func Ping(data []byte) ([]byte, error) {
	return []byte(SignPingRET), nil
}
