import { userV1ServiceBuilder } from '../../userV1ServiceBuilder';
import {
  userV1InfoInputParameterSchema,
  userV1InfoInputPayloadSchema,
  userV1InfoOutputPayloadSchema,
} from './schema';

export const infoCommandBuilder = userV1ServiceBuilder
  .getCommandBuilder('info', 'Returns info about user')
  .addPayloadSchema(userV1InfoInputPayloadSchema)
  .addParameterSchema(userV1InfoInputParameterSchema)
  .addOutputSchema(userV1InfoOutputPayloadSchema)
  .exposeAsHttpEndpoint('GET', 'user/info')
  .addOpenApiTags('User')
  .setCommandFunction(async function (_context, _payload, _parameter) {
    // add your business logic here
  });
