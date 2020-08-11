import grpc

from proto.sample_service import sample_service_pb2, sample_service_pb2_grpc


def multiply(a, b):
    # connect rpc server
    channel = grpc.insecure_channel('localhost:50052')
    # invoke rpc method
    stub = sample_service_pb2_grpc.SampleStub(channel)
    response = stub.Multiply(sample_service_pb2.MultiplyReq(a=a, b=b))
    print("sample_service.Multiple result: \n{}".format(response.result))


def get_user_info(uid):
    # connect rpc server
    channel = grpc.insecure_channel('localhost:50052')
    # invoke rpc method
    stub = sample_service_pb2_grpc.SampleStub(channel)
    try:
        response = stub.GetUserById(sample_service_pb2.GetUserByIdReq(uid=uid))
        print("sample_service.GetUserById result: \n{}".format(response))
    except grpc.RpcError as e:
        print("sample_service.GetUserById got error: \n{}".format(e))


def create_user_from_external(uid):
    # connect rpc server
    channel = grpc.insecure_channel('localhost:50052')
    # invoke rpc method
    stub = sample_service_pb2_grpc.SampleStub(channel)
    try:
        response = stub.CreateUserFromExternal(sample_service_pb2.CreateUserFromExternalReq(externalUid=uid))
        print("sample_service.CreateUserFromExternal returned uid: {}".format(response.uid))
    except grpc.RpcError as e:
        print("sample_service.CreateUserFromExternal got error: \n{}".format(e))


if __name__ == '__main__':
    multiply(2, 3)

    # get_user_info(2)

    # uids = [0, 1, 2, 3, 4, 5]
    # for one in uids:
    #     print('external uid', one)
    #     create_user_from_external(one)
