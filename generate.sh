#!/bin/bash

protoc greet-app/greetpb/greet.proto --go_out=plugins=grpc:.

protoc calculator-app/calculatorpb/calculator.proto --go_out=plugins=grpc:.

protoc blog-app/blogpb/blog.proto --go_out=plugins=grpc:.
