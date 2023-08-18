//
//  Generated code. Do not modify.
//  source: user/user.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import 'user.pb.dart' as $0;

export 'user.pb.dart';

@$pb.GrpcServiceName('user.UserService')
class UserServiceClient extends $grpc.Client {
  static final _$say = $grpc.ClientMethod<$0.SayRequest, $0.SayResponse>(
      '/user.UserService/Say',
      ($0.SayRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.SayResponse.fromBuffer(value));

  UserServiceClient($grpc.ClientChannel channel,
      {$grpc.CallOptions? options,
      $core.Iterable<$grpc.ClientInterceptor>? interceptors})
      : super(channel, options: options,
        interceptors: interceptors);

  $grpc.ResponseFuture<$0.SayResponse> say($0.SayRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$say, request, options: options);
  }
}

@$pb.GrpcServiceName('user.UserService')
abstract class UserServiceBase extends $grpc.Service {
  $core.String get $name => 'user.UserService';

  UserServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.SayRequest, $0.SayResponse>(
        'Say',
        say_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.SayRequest.fromBuffer(value),
        ($0.SayResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.SayResponse> say_Pre($grpc.ServiceCall call, $async.Future<$0.SayRequest> request) async {
    return say(call, await request);
  }

  $async.Future<$0.SayResponse> say($grpc.ServiceCall call, $0.SayRequest request);
}
