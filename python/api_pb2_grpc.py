# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

import api_pb2 as api__pb2


class WelcomeStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.Hi = channel.unary_stream(
        '/api.Welcome/Hi',
        request_serializer=api__pb2.Hello.SerializeToString,
        response_deserializer=api__pb2.Status.FromString,
        )
    self.Translate = channel.stream_stream(
        '/api.Welcome/Translate',
        request_serializer=api__pb2.ToTranslate.SerializeToString,
        response_deserializer=api__pb2.TranslateResult.FromString,
        )
    self.Store = channel.stream_unary(
        '/api.Welcome/Store',
        request_serializer=api__pb2.RecordMessage.SerializeToString,
        response_deserializer=api__pb2.StoreSummary.FromString,
        )


class WelcomeServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def Hi(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Translate(self, request_iterator, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Store(self, request_iterator, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_WelcomeServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'Hi': grpc.unary_stream_rpc_method_handler(
          servicer.Hi,
          request_deserializer=api__pb2.Hello.FromString,
          response_serializer=api__pb2.Status.SerializeToString,
      ),
      'Translate': grpc.stream_stream_rpc_method_handler(
          servicer.Translate,
          request_deserializer=api__pb2.ToTranslate.FromString,
          response_serializer=api__pb2.TranslateResult.SerializeToString,
      ),
      'Store': grpc.stream_unary_rpc_method_handler(
          servicer.Store,
          request_deserializer=api__pb2.RecordMessage.FromString,
          response_serializer=api__pb2.StoreSummary.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'api.Welcome', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
