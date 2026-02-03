/**
 * @fileoverview gRPC-Web generated client stub for event
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.event = require('./eventservice_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.event.EventServiceClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.event.EventServicePromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.event.Events,
 *   !proto.event.Response>}
 */
const methodDescriptor_EventService_dispatch = new grpc.web.MethodDescriptor(
  '/event.EventService/dispatch',
  grpc.web.MethodType.UNARY,
  proto.event.Events,
  proto.event.Response,
  /**
   * @param {!proto.event.Events} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.event.Response.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.event.Events,
 *   !proto.event.Response>}
 */
const methodInfo_EventService_dispatch = new grpc.web.AbstractClientBase.MethodInfo(
  proto.event.Response,
  /**
   * @param {!proto.event.Events} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.event.Response.deserializeBinary
);


/**
 * @param {!proto.event.Events} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.event.Response)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.event.Response>|undefined}
 *     The XHR Node Readable Stream
 */
proto.event.EventServiceClient.prototype.dispatch =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/event.EventService/dispatch',
      request,
      metadata || {},
      methodDescriptor_EventService_dispatch,
      callback);
};


/**
 * @param {!proto.event.Events} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.event.Response>}
 *     Promise that resolves to the response
 */
proto.event.EventServicePromiseClient.prototype.dispatch =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/event.EventService/dispatch',
      request,
      metadata || {},
      methodDescriptor_EventService_dispatch);
};


module.exports = proto.event;

