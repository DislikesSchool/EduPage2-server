import {
  getCommandContextMock,
  getEventBridgeMock,
  getLoggerMock,
} from '@purista/core';
import { createSandbox } from 'sinon';

import { icanteenV1Service } from '../../icanteenV1Service';
import { lunchesCommandBuilder } from './lunchesCommandBuilder';
import {
  IcanteenV1LunchesInputParameter,
  IcanteenV1LunchesInputPayload,
} from './types';

describe('service Icanteen version 1 - command lunches', () => {
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

    const lunches = lunchesCommandBuilder.getCommandFunction().bind(service);

    const payload: IcanteenV1LunchesInputPayload = undefined;

    const parameter: IcanteenV1LunchesInputParameter = {};

    const context = getCommandContextMock(payload, parameter, sandbox);

    const result = await lunches(context.mock, payload, parameter);

    expect(result).toBeUndefined();
  });
});
