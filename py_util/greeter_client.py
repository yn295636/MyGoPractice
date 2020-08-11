import grpc

from proto.greeter_service import greeter_service_pb2, greeter_service_pb2_grpc


def save_and_get_from_redis():
    # connect rpc server
    channel = grpc.insecure_channel('localhost:50051')
    # invoke rpc method
    stub = greeter_service_pb2_grpc.GreeterStub(channel)

    data = greeter_service_pb2.StoreInRedisRequest(
        key="my_k",
        value="my_value"
    )
    resp = stub.StoreInRedis(data)
    print("Save {{{}}} into redis result: {}".format(data, resp.result))

    resp = stub.GetFromRedis(greeter_service_pb2.GetFromRedisRequest(key=data.key))
    print("Get key={} value from redis: {}".format(data.key, resp.value))


if __name__ == '__main__':
    save_and_get_from_redis()
