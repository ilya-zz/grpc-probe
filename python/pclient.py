import api_pb2_grpc as stub
import api_pb2 as pb
import grpc

channel = grpc.insecure_channel("localhost:7777")
client = stub.WelcomeStub(channel)

def inputIt():
    d = [
        "foo foo foo bar",
            "foo goo foo bar"]
    id = 1
    for s in d:        
        id += 1
        yield pb.ToTranslate( id = id, text = s)

data = [
    pb.ToTranslate( id = 1, text = "zhopa zhopa"),
    pb.ToTranslate( id = 2, text = "qqq qqq qqq"),
]


for msg in client.Translate(iter(data)):
    print(msg)
