import 'dotenv/config';

import { fastifyFormbody } from '@fastify/formbody';
import { AmqpBridge, AmqpBridgeConfig } from '@purista/amqpbridge';
import { httpServerV1Service } from '@purista/httpserver';

import httpServerConfig from '../config/httpServerConfig';
import { edupageV1Service } from './service/edupage/v1/edupageV1Service';
import { icanteenV1Service } from './service/icanteen/v1/icanteenV1Service';
import { userV1Service } from './service/user/v1/userV1Service';

export const main = async () => {
  // initiate the event bridge as first step
  if (process.env.AMQP_URL === undefined) {
    return;
  }
  const config: AmqpBridgeConfig = {
    url: process.env.AMQP_URL,
  };

  const eventBridge = new AmqpBridge(config);
  await eventBridge.start();
  // initiate the webserver service as second step
  const httpServerService = httpServerV1Service.getInstance(eventBridge, {
    serviceConfig: httpServerConfig,
  });
  httpServerService.server?.register(fastifyFormbody);

  const userInstance = userV1Service.getInstance(eventBridge);
  const icanteenInstance = icanteenV1Service.getInstance(eventBridge);
  const edupageInstance = edupageV1Service.getInstance(eventBridge);
  // initiate/start the user instance
  // it registers the commands and the subscriptions to the event bridge
  await userInstance.start();
  await icanteenInstance.start();
  await edupageInstance.start();

  // start the webserver
  await httpServerService.start();

  // add initiation and start of services here
};

main();
