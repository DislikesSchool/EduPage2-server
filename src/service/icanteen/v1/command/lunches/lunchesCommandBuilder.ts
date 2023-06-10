import { icanteenV1ServiceBuilder } from '../../icanteenV1ServiceBuilder';
import {
  icanteenV1LunchesInputParameterSchema,
  icanteenV1LunchesInputPayloadSchema,
  icanteenV1LunchesOutputPayloadSchema,
} from './schema';

export const lunchesCommandBuilder = icanteenV1ServiceBuilder
  .getCommandBuilder(
    'lunches',
    'Returns list of lunch options for the following month',
  )
  .addPayloadSchema(icanteenV1LunchesInputPayloadSchema)
  .addParameterSchema(icanteenV1LunchesInputParameterSchema)
  .addOutputSchema(icanteenV1LunchesOutputPayloadSchema)
  .exposeAsHttpEndpoint('GET', 'icanteen/lunches')
  .addOpenApiTags('Lunches')
  .setCommandFunction(async function (_context, _payload, _parameter) {
    // add your business logic here
  });
