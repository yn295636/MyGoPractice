import json

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


def save_and_query_from_mongo():
    # connect rpc server
    channel = grpc.insecure_channel('localhost:50051')
    # invoke rpc method
    stub = greeter_service_pb2_grpc.GreeterStub(channel)

    data1 = greeter_service_pb2.StoreInMongoRequest(data=json.dumps(
        {
            'my_name': '111',
            'my_age': 18
        }
    ))
    data2 = greeter_service_pb2.StoreInMongoRequest(data=json.dumps(
        {
            'my_name': '222',
            'my_age': 20
        }
    ))
    resp = stub.StoreInMongo(data1)
    print("Save {{{}}} into mongo result: {}".format(data1, resp.result))
    resp = stub.StoreInMongo(data2)
    print("Save {{{}}} into mongo result: {}".format(data2, resp.result))

    resp = stub.CountOfDataInMongo(greeter_service_pb2.CountOfDataInMongoRequest())
    print("Count of data in mongo: {}".format(resp.count))


if __name__ == '__main__':
    save_and_query_from_mongo()
    # save_and_get_from_redis()
