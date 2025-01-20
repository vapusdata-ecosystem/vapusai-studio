from server.grpc import GrpcServer

if __name__ == "__main__":
    GrpcServer().start()

# import grpc
# from grpc_reflection.v1alpha import reflection
# from concurrent import futures
# from controller import vapusmlutilities
# from utils.importer import proto_importer
# proto_importer()

# from protos.vapus_aiutilities.v1alpha1 import vapus_aiutilities_pb2_grpc
# import protos.vapus_aiutilities.v1alpha1.vapus_aiutilities_pb2 as pb2


# def serve():
#     # Create a gRPC server
#     server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    
#     # Add the service to the server
#     vapus_aiutilities_pb2_grpc.add_AIUtilityServicer_to_server(vapusmlutilities.AIUtilityServicer(), server)
    
#     SERVICE_NAMES = (
#         pb2.DESCRIPTOR.services_by_name['AIUtility'].full_name,
#         reflection.SERVICE_NAME,
#     )
#     reflection.enable_server_reflection(SERVICE_NAMES, server)

#     # Bind the server to a specific port
#     server.add_insecure_port('[::]:50051')
    
#     # Start the server
#     print("gRPC server running on port 50051")
#     server.start()
    
#     # Keep the server running
#     server.wait_for_termination()

# if __name__ == "__main__":
#     serve()
