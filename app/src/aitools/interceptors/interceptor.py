import grpc
from typing import Callable, Any
from helpers.jwtvalidate import validateJWT
import json
from grpc_interceptor import ServerInterceptor
from grpc_interceptor.exceptions import GrpcException

class Interceptor(ServerInterceptor):
    """
    A gRPC interceptor that processes requests through multiple functions before handling the actual RPC.
    """
    
    def intercept(self, method: Callable, request: Any, context: grpc.ServicerContext, method_name: str) -> Any:
        try:
            
            metadata = context.invocation_metadata()
            jwt = self.extract_jwt(metadata)
            scope = validateJWT(jwt)
            context.scope = scope  
        except grpc.RpcError as e:
            context.abort(e.code(), e.details())
        except Exception:
            context.abort(grpc.StatusCode.INTERNAL, "Unexpected server error")
        
        return method(request, context)

    def extract_jwt(self, metadata):
        """
        Extracts the JWT token from the metadata.
        """
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
