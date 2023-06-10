import { ServiceEvent } from '../../../../ServiceEvent.enum';
import { userV1ServiceBuilder } from '../../userV1ServiceBuilder';
import {
  userV1RegisterInputPayloadSchema,
  userV1RegisterOutputPayloadSchema,
} from './schema';

export const registerCommandBuilder = userV1ServiceBuilder
  .getCommandBuilder(
    'register',
    'Registers a new user by saving their EduPage credentials',
  )
  .setSuccessEventName(ServiceEvent.UserRegistered)
  .addPayloadSchema(userV1RegisterInputPayloadSchema)
  .addOpenApiTags('User')
  .addOutputSchema(userV1RegisterOutputPayloadSchema)
  .exposeAsHttpEndpoint('POST', 'user/register')
  .setCommandFunction(async function (_context, _payload, _parameter) {
    this.logger.info(_context);
    this.logger.info(_payload);
    return {
      token: 'hi',
    };
  });
