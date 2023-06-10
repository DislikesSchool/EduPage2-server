import {
  getCommandContextMock,
  getEventBridgeMock,
  getLoggerMock,
} from '@purista/core';
import { createSandbox } from 'sinon';

import { icanteenV1Service } from '../../icanteenV1Service';
import { setupCommandBuilder } from './setupCommandBuilder';
import {
  IcanteenV1SetupInputParameter,
  IcanteenV1SetupInputPayload,
} from './types';

describe('service Icanteen version 1 - command setup', () => {
  let sandbox = createSandbox();
  beforeEach(() => {
    sandbox = createSandbox();
  });

  afterEach(() => {
    sandbox.restore();
  });

  test('does not throw', async () => {
    const service = icanteenV1Service.getInstance(
      getEventBridgeMock(sandbox).mock,
      { logger: getLoggerMock(sandbox).mock },
    );

    const setup = setupCommandBuilder.getCommandFunction().bind(service);

    const payload: IcanteenV1SetupInputPayload = undefined;

    const parameter: IcanteenV1SetupInputParameter = {};

    const context = getCommandContextMock(payload, parameter, sandbox);

    const result = await setup(context.mock, payload, parameter);

    expect(result).toBeUndefined();
  });
});
