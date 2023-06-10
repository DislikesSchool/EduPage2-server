import { userV1ServiceBuilder } from '../../userV1ServiceBuilder';
import {
  userV1ValidateInputParameterSchema,
  userV1ValidateInputPayloadSchema,
  userV1ValidateOutputPayloadSchema,
} from './schema';

export const validateCommandBuilder = userV1ServiceBuilder
  .getCommandBuilder('validate', 'Validates a access token')
  .addPayloadSchema(userV1ValidateInputPayloadSchema)
  .addParameterSchema(userV1ValidateInputParameterSchema)
  .addOutputSchema(userV1ValidateOutputPayloadSchema)
  .exposeAsHttpEndpoint('GET', 'user/validate')
  .addOpenApiTags('User')
  .setCommandFunction(async function (_context, _payload, _parameter) {
    return {
      valid: true,
    };
  });
