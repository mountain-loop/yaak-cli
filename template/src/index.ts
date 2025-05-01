import {PluginDefinition} from "@yaakapp/api";

export const plugin: PluginDefinition = {
    httpRequestActions: [
        {
            label: "Hello, From Plugin",
            icon: "info",
            async onSelect(ctx, args) {
                await ctx.toast.show({
                    color: "success",
                    message: `You clicked the request ${args.httpRequest.id}`
                });
            },
        },
    ],
};
