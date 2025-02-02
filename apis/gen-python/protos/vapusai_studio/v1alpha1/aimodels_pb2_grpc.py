# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from protos.vapusai_studio.v1alpha1 import aimodels_pb2 as protos_dot_vapusai__studio_dot_v1alpha1_dot_aimodels__pb2


class AIModelsStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Manager = channel.unary_unary(
                '/vapusai.v1alpha1.AIModels/Manager',
                request_serializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aimodels__pb2.AIModelNodeConfiguratorRequest.SerializeToString,
                response_deserializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aimodels__pb2.AIModelNodeResponse.FromString,
                _registered_method=True)
        self.Getter = channel.unary_unary(
                '/vapusai.v1alpha1.AIModels/Getter',
                request_serializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aimodels__pb2.AIModelNodeGetterRequest.SerializeToString,
                response_deserializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aimodels__pb2.AIModelNodeResponse.FromString,
                _registered_method=True)


class AIModelsServicer(object):
    """Missing associated documentation comment in .proto file."""

    def Manager(self, request, context):
        """*
        Retrieves a data worker deployment.
        @param {Manager} request - The request object containing the AI model node configurator.
        @returns {AIModelNodeResponse} The response object containing the AI model node configurator.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Getter(self, request, context):
        """*
        Retrieves a data worker deployment.
        @param {AIModelNodeGetterRequest} request - The request object containing the AI model node getter.
        @returns {AIModelNodeResponse} The response object containing the AI model node getter.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_AIModelsServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Manager': grpc.unary_unary_rpc_method_handler(
                    servicer.Manager,
                    request_deserializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aimodels__pb2.AIModelNodeConfiguratorRequest.FromString,
                    response_serializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aimodels__pb2.AIModelNodeResponse.SerializeToString,
            ),
            'Getter': grpc.unary_unary_rpc_method_handler(
                    servicer.Getter,
                    request_deserializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aimodels__pb2.AIModelNodeGetterRequest.FromString,
                    response_serializer=protos_dot_vapusai__studio_dot_v1alpha1_dot_aimodels__pb2.AIModelNodeResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'vapusai.v1alpha1.AIModels', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('vapusai.v1alpha1.AIModels', rpc_method_handlers)


 # This class is part of an EXPERIMENTAL API.
class AIModels(object):
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
            '/vapusai.v1alpha1.AIModels/Manager',
            protos_dot_vapusai__studio_dot_v1alpha1_dot_aimodels__pb2.AIModelNodeConfiguratorRequest.SerializeToString,
            protos_dot_vapusai__studio_dot_v1alpha1_dot_aimodels__pb2.AIModelNodeResponse.FromString,
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
            '/vapusai.v1alpha1.AIModels/Getter',
            protos_dot_vapusai__studio_dot_v1alpha1_dot_aimodels__pb2.AIModelNodeGetterRequest.SerializeToString,
            protos_dot_vapusai__studio_dot_v1alpha1_dot_aimodels__pb2.AIModelNodeResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)
