import type { NuxtAxiosInstance } from '@nuxtjs/axios';
import { z } from 'zod';

import { makeAPI } from './makeAPI';
import { LikeCoinMigrationSchema } from './models/likeCoinMigration';

const ResponseSchema = z.object({
  migration: LikeCoinMigrationSchema,
});

export const makeCancelLikeCoinMigrationAPI =
  (axiosInstance: NuxtAxiosInstance) => (cosmosAddress: string) =>
    makeAPI({
      method: 'DELETE',
      url: `/likecoin/migration/${cosmosAddress}`,
      responseSchema: ResponseSchema,
    })(axiosInstance)();

export type CancelLikeCoinMigration = ReturnType<
  typeof makeCancelLikeCoinMigrationAPI
>;
