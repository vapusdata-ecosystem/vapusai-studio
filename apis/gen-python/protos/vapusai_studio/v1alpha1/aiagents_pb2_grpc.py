# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from protos.vapusai_studio.v1alpha1 import aiagents_pb2 as protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2


class AIAgentsStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Manager = channel.unary_unary(
                '/vapusai.v1alpha1.AIAgents/Manager',
                request_serializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentManagerRequest.SerializeToString,
                response_deserializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentResponse.FromString,
                _registered_method=True)
        self.Getter = channel.unary_unary(
                '/vapusai.v1alpha1.AIAgents/Getter',
                request_serializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentGetterRequest.SerializeToString,
                response_deserializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentResponse.FromString,
                _registered_method=True)


class AIAgentsServicer(object):
    """Missing associated documentation comment in .proto file."""

    def Manager(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Getter(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_AIAgentsServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Manager': grpc.unary_unary_rpc_method_handler(
                    servicer.Manager,
                    request_deserializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentManagerRequest.FromString,
                    response_serializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentResponse.SerializeToString,
            ),
            'Getter': grpc.unary_unary_rpc_method_handler(
                    servicer.Getter,
                    request_deserializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentGetterRequest.FromString,
                    response_serializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'vapusai.v1alpha1.AIAgents', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('vapusai.v1alpha1.AIAgents', rpc_method_handlers)


 # This class is part of an EXPERIMENTAL API.
class AIAgents(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def Manager(request,
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
            '/vapusai.v1alpha1.AIAgents/Manager',
            protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentManagerRequest.SerializeToString,
            protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentResponse.FromString,
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
    def Getter(request,
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
            '/vapusai.v1alpha1.AIAgents/Getter',
            protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentGetterRequest.SerializeToString,
            protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)


class AIAgentStudioStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.ChatStream = channel.unary_stream(
                '/vapusai.v1alpha1.AIAgentStudio/ChatStream',
                request_serializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentInvokeRequest.SerializeToString,
                response_deserializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentInvokeStreamResponse.FromString,
                _registered_method=True)
        self.Chat = channel.unary_unary(
                '/vapusai.v1alpha1.AIAgentStudio/Chat',
                request_serializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentInvokeRequest.SerializeToString,
                response_deserializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentInvokeResponse.FromString,
                _registered_method=True)


class AIAgentStudioServicer(object):
    """Missing associated documentation comment in .proto file."""

    def ChatStream(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Chat(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_AIAgentStudioServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'ChatStream': grpc.unary_stream_rpc_method_handler(
                    servicer.ChatStream,
                    request_deserializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentInvokeRequest.FromString,
                    response_serializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentInvokeStreamResponse.SerializeToString,
            ),
            'Chat': grpc.unary_unary_rpc_method_handler(
                    servicer.Chat,
                    request_deserializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentInvokeRequest.FromString,
                    response_serializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentInvokeResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'vapusai.v1alpha1.AIAgentStudio', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('vapusai.v1alpha1.AIAgentStudio', rpc_method_handlers)


 # This class is part of an EXPERIMENTAL API.
class AIAgentStudio(object):
    """Missing associated documentation comment in .proto file."""

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
            '/vapusai.v1alpha1.AIAgentStudio/ChatStream',
            protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentInvokeRequest.SerializeToString,
            protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentInvokeStreamResponse.FromString,
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
            '/vapusai.v1alpha1.AIAgentStudio/Chat',
            protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentInvokeRequest.SerializeToString,
            protos_dot_vapusai__studio_dot_v1alpha1_dot_aiagents__pb2.AgentInvokeResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)
