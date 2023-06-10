import { ServiceEvent } from '../../../../ServiceEvent.enum';
import { icanteenV1ServiceBuilder } from '../../icanteenV1ServiceBuilder';
import {
  icanteenV1SetupInputParameterSchema,
  icanteenV1SetupInputPayloadSchema,
  icanteenV1SetupOutputPayloadSchema,
} from './schema';

export const setupCommandBuilder = icanteenV1ServiceBuilder
  .getCommandBuilder('setup', 'Sets up the iCanteen integration for a user')
  .setSuccessEventName(ServiceEvent.IcanteenSetup)
  .addPayloadSchema(icanteenV1SetupInputPayloadSchema)
  .addParameterSchema(icanteenV1SetupInputParameterSchema)
  .addOutputSchema(icanteenV1SetupOutputPayloadSchema)
  .exposeAsHttpEndpoint('POST', 'icanteen/setup')
  .addOpenApiTags('Lunches')
  .setCommandFunction(async function (_context, _payload, _parameter) {
    // add your business logic here
  });
