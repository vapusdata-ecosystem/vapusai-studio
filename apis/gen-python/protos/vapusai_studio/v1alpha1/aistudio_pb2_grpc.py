# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from protos.vapusai_studio.v1alpha1 import aistudio_pb2 as protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2


class AIModelStudioStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.GenerateEmbeddings = channel.unary_unary(
                '/vapusai.v1alpha1.AIModelStudio/GenerateEmbeddings',
                request_serializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2.EmbeddingsInterface.SerializeToString,
                response_deserializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2.EmbeddingsResponse.FromString,
                _registered_method=True)
        self.Chat = channel.unary_unary(
                '/vapusai.v1alpha1.AIModelStudio/Chat',
                request_serializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2.ChatRequest.SerializeToString,
                response_deserializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2.ChatResponse.FromString,
                _registered_method=True)
        self.ChatStream = channel.unary_stream(
                '/vapusai.v1alpha1.AIModelStudio/ChatStream',
                request_serializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2.ChatRequest.SerializeToString,
                response_deserializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2.ChatResponse.FromString,
                _registered_method=True)


class AIModelStudioServicer(object):
    """Missing associated documentation comment in .proto file."""

    def GenerateEmbeddings(self, request, context):
        """Generates embeddings for the given input text.
        @param {EmbeddingsInterface} request - The request object containing the embeddings interface.
        @returns {EmbeddingsResponse} The response object containing the embeddings.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Chat(self, request, context):
        """Generates content based on the given prompt parameters.
        @param {ChatRequest} request - The request object containing the generate interface.
        @returns {ChatResponse} The response object containing the generated content.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ChatStream(self, request, context):
        """Generates content in a streaming fashion based on the given prompt parameters.
        @param {ChatRequest} request - The request object containing the generate interface.
        @returns {stream GenerateStreamResponse} The response object containing the generated content in a stream.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_AIModelStudioServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'GenerateEmbeddings': grpc.unary_unary_rpc_method_handler(
                    servicer.GenerateEmbeddings,
                    request_deserializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2.EmbeddingsInterface.FromString,
                    response_serializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2.EmbeddingsResponse.SerializeToString,
            ),
            'Chat': grpc.unary_unary_rpc_method_handler(
                    servicer.Chat,
                    request_deserializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2.ChatRequest.FromString,
                    response_serializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2.ChatResponse.SerializeToString,
            ),
            'ChatStream': grpc.unary_stream_rpc_method_handler(
                    servicer.ChatStream,
                    request_deserializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2.ChatRequest.FromString,
                    response_serializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2.ChatResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'vapusai.v1alpha1.AIModelStudio', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('vapusai.v1alpha1.AIModelStudio', rpc_method_handlers)


 # This class is part of an EXPERIMENTAL API.
class AIModelStudio(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def GenerateEmbeddings(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/vapusai.v1alpha1.AIModelStudio/GenerateEmbeddings',
            protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2.EmbeddingsInterface.SerializeToString,
            protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2.EmbeddingsResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def Chat(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/vapusai.v1alpha1.AIModelStudio/Chat',
            protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2.ChatRequest.SerializeToString,
            protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2.ChatResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def ChatStream(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_stream(
            request,
            target,
            '/vapusai.v1alpha1.AIModelStudio/ChatStream',
            protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2.ChatRequest.SerializeToString,
            protos_dot_vapusai__studio_dot_v1alpha1_dot_aistudio__pb2.ChatResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)
