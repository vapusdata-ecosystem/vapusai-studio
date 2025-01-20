from utils.importer import proto_importer
from presidio_analyzer import AnalyzerEngine
from presidio_anonymizer import AnonymizerEngine
from presidio_anonymizer.entities import RecognizerResult, OperatorConfig
from google.protobuf.json_format import MessageToDict
import json

proto_importer()

import grpc
from grpc import ServerInterceptor
from protos.vapus_aiutilities.v1alpha1 import vapus_aiutilities_pb2_grpc as aiutilities
import protos.vapus_aiutilities.v1alpha1.vapus_aiutilities_pb2 as pb2
from sentence_transformers import SentenceTransformer
from middlewares import jwtvalidate
from models import aiutility_models
from enum import Enum

class AIUtilityServicer(aiutilities.AIUtilityServicer):

    def getJwt(self,context):
        metadata =  context.invocation_metadata()
        token = None
        for key, value in metadata:
            if key.lower() == "authorization":
                if value.startswith("Bearer "):
                    token = value[len("Bearer "):]  
                else:
                    context.set_code(grpc.StatusCode.UNAUTHENTICATED)
                    context.set_details("Invalid Authorization header format")
                    return None 
        return token
    
    def GenerateEmbedding(self, request, context):
        
        
        try:
            aiutility_models.GenerateEmbeddingRequest.model_validate(MessageToDict(request))
        except Exception as e:
            raise e
        
        
        jwt = self.getJwt(context)
        
        if jwt is None:
            context.set_code(grpc.StatusCode.UNAUTHENTICATED)
            context.set_details("Authorization header missing")
            return pb2.GenerateEmbeddingResponse()

        try:
            jwtvalidate.validateJWT(jwt)
        except Exception as e:
            raise e
        
        try:
           
            model = SentenceTransformer("all-MiniLM-L6-v2")

            sentences = request.text 
       
            embeddings = model.encode(sentences)

            response = pb2.GenerateEmbeddingResponse()

           
            for idx, embedding in enumerate(embeddings):
                embedding_proto = response.Embeddings()
                embedding_proto.embedding.extend(embedding.tolist())  
                embedding_proto.index = idx
                response.embeddings.append(embedding_proto)

            return response

        except Exception as e:
          
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"An error occurred: {str(e)}")
            return pb2.GenerateEmbeddingResponse()
        
    def SensitivityAnalyzer(self, request, context):

        try:
            aiutility_models.SensitivityAnalyzerRequest.model_validate(MessageToDict(request))
        except Exception as e:
            raise e
        
        jwt = self.getJwt(context)
        
        if jwt is None:
            context.set_code(grpc.StatusCode.UNAUTHENTICATED)
            context.set_details("Authorization header missing")
            return pb2.SensitivityAnalyzerResponse()
        
        try:
            jwtvalidate.validateJWT(jwt)
        except Exception as e:
            raise e
        
        
        analyzer = AnalyzerEngine()
        
        response = pb2.SensitivityAnalyzerResponse()

        
        if request.action == 0:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details("Invalid action specified")
            return response

        
        if request.action == 1:
            '''
                analyze
            '''
            
            for index,text in enumerate(request.text):
               
                try:
                    result = analyzer.analyze(text = text ,language="en",entities=request.entities)
                except Exception as e:
                    raise e
                
                
                processedOutput = response.ProcessedOutput()
                processedOutput.text = text
                processedOutput.index = index
                for entity in request.entities:
                    processedOutput.entities.append(entity)
                processedOutput.action  = request.postDetectAction
                
                for res in result:

                    item_dict = res.to_dict()                   
                    analyzedOutput = pb2.AnalyzedOutput()
                    analyzedOutput.type = item_dict.get("entity_type")
                    analyzedOutput.start = item_dict.get("start")
                    analyzedOutput.end  = item_dict.get("end")
                    analyzedOutput.score = item_dict.get("score")
                    
                    processedOutput.AnalyzedOutputs.append(analyzedOutput)
                    
                response.output.append(processedOutput)
            return response
        if request.action == 2:
            '''
                act
            '''

            if request.postDetectAction == 0:
                context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
                context.set_details("Invalid action specified")
                return response

            for index,text in enumerate(request.text):
                result = analyzer.analyze(text = text ,language="en",entities=request.entities)
                processedOutput = response.ProcessedOutput()
                processedOutput.text = text
                processedOutput.index = index
                for entity in request.entities:
                    processedOutput.entities.append(entity)
                processedOutput.action  = request.postDetectAction
                analyzer_results = []

                
                for res in result:

                    item_dict = res.to_dict()                   
                    analyzedOutput = pb2.AnalyzedOutput()
                    analyzedOutput.type = item_dict.get("entity_type")
                    analyzedOutput.start = item_dict.get("start")
                    analyzedOutput.end  = item_dict.get("end")
                    analyzedOutput.score = item_dict.get("score")
                    
                    
                    if item_dict.get("score")>=0.7:
                        analyzer_results.append(res)
                    
                    processedOutput.AnalyzedOutputs.append(analyzedOutput)
                operators = {}
                
                placeholder = {1:"xxxx",2:"Placeholder",3:""}

               
                for item in analyzer_results:
                    operators[item.entity_type] = OperatorConfig("replace", {"new_value": placeholder.get(request.postDetectAction)})
                
                engine = AnonymizerEngine()
                editedText = engine.anonymize(text = text,analyzer_results=analyzer_results,operators=operators)
                processedOutput.text = editedText.text
                response.output.append(processedOutput)
            return response





        








