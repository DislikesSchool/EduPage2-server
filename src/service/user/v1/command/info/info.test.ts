import {
  getCommandContextMock,
  getEventBridgeMock,
  getLoggerMock,
} from '@purista/core';
import { createSandbox } from 'sinon';

import { userV1Service } from '../../userV1Service';
import { infoCommandBuilder } from './infoCommandBuilder';
import { UserV1InfoInputParameter, UserV1InfoInputPayload } from './types';

describe('service User version 1 - command info', () => {
  let sandbox = createSandbox();
  beforeEach(() => {
    sandbox = createSandbox();
  });

  afterEach(() => {
    sandbox.restore();
  });

  test('does not throw', async () => {
    const service = userV1Service.getInstance(
      getEventBridgeMock(sandbox).mock,
      { logger: getLoggerMock(sandbox).mock },
    );

    const info = infoCommandBuilder.getCommandFunction().bind(service);

    const payload: UserV1InfoInputPayload = undefined;

    const parameter: UserV1InfoInputParameter = {};

    const context = getCommandContextMock(payload, parameter, sandbox);

    const result = await info(context.mock, payload, parameter);

    expect(result).toBeUndefined();
  });
});
