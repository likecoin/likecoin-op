import { z } from 'zod';

export const ConfigSchema = z.object({
  isTestnet: z.boolean().catch(true),
  authcoreRedirectUrl: z.string(),
  apiBaseURL: z.string(),
  likerlandUrlBase: z.string(),
  crispWebsiteId: z.string(),
});

export type Config = z.infer<typeof ConfigSchema>;
