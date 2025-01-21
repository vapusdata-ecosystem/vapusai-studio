import grpc
from typing import Callable,Any
from helpers.jwtvalidate import validateJWT
import json
from grpc_interceptor import ServerInterceptor
from grpc_interceptor.exceptions import GrpcException


class Interceptor(grpc.ServerInterceptor):
    """
    A gRPC interceptor that processes requests through multiple functions before handling the actual RPC.
    """
    
    def intercept_service(
        self,
        continuation: Callable[[grpc.HandlerCallDetails], grpc.RpcMethodHandler],
        handler_call_details: grpc.HandlerCallDetails,
    ) -> grpc.RpcMethodHandler:
        """
        Intercepts incoming gRPC calls and processes them through multiple functions.
        """

        print(f"Intercepting call to method: {handler_call_details.method}")

        
        try:
            handler_call_details = self.validate_jwt(handler_call_details)
        except grpc.RpcError as error:
            raise error

        
        return continuation(handler_call_details)

    def extract_jwt(self, metadata):
       
        for key, value in metadata:
            if key.lower() == "authorization":
                if value.startswith("Bearer "):
                    return value[len("Bearer "):]  
                else:
                    raise grpc.RpcError(
                        grpc.StatusCode.UNAUTHENTICATED,
                        "Invalid Authorization header format",
                    )
        raise grpc.RpcError(
            grpc.StatusCode.UNAUTHENTICATED,
            "Authorization metadata is missing",
        )

    def validate_jwt(self, handler_call_details: grpc.HandlerCallDetails) -> grpc.HandlerCallDetails:
        
        
        metadata = dict(handler_call_details.invocation_metadata)
        
       
        jwt = self.extract_jwt(handler_call_details.invocation_metadata)
        
        try:
            scope = validateJWT(jwt)
        except Exception as e:
            raise grpc.RpcError(
                grpc.StatusCode.PERMISSION_DENIED,
                "Invalid JWT",
            )

        
        metadata["scope"] = json.dumps(scope)
        new_metadata = tuple(metadata.items())
        
        '''
            need to debug adding scope to metadata
        '''
        # try:
        #     updated_call_details = grpc.HandlerCallDetails(
        #         method=handler_call_details.method,
        #         invocation_metadata=new_metadata,
        #     )
        # except Exception as e:
        #     print(e)
        #     raise e
        
        return handler_call_details
