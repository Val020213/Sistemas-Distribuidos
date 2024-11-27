"use server";

import { revalidatePath, revalidateTag } from "next/cache";

export const revalidateServerPath = async (path: string) => {
  revalidatePath(path);
};

export const revalidateServerTags = async (tag: string) => {
  revalidateTag(tag);
};
