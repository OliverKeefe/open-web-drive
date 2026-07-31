import type { KcContext } from "../KcContext";
import { useI18n } from "../i18n";
import { cn } from "@/lib/utils.ts";
import { Button } from "@/components/ui/button.tsx";
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card.tsx";
import {
    Field,
    FieldGroup,
    FieldLabel,
    FieldDescription,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input.tsx";

export default function Login(props: { kcContext: Extract<KcContext, { pageId: "login.ftl" }> }) {
    const { kcContext } = props;
    const { url, realm, auth, messagesPerField, social } = kcContext;

    const { i18n } = useI18n({ kcContext });

    const usernameError = messagesPerField.existsError("username", "password");

    return (
        <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
            <div className="w-full max-w-sm">
                <Card>
                    <CardHeader>
                        <CardTitle>{i18n.msg("loginTitle", kcContext.realm.displayName ?? kcContext.realm.name)}</CardTitle>
                        <CardDescription>
                            Sign in with your OWD account.
                        </CardDescription>
                    </CardHeader>

                    <CardContent>
                        <form
                            id="kc-form-login"
                            action={url.loginAction}
                            method="post"
                        >
                            <FieldGroup>
                                {/* Username */}
                                <Field>
                                    <FieldLabel htmlFor="username">
                                        {realm.loginWithEmailAllowed
                                            ? i18n.msg("email")
                                            : i18n.msg("username")}
                                    </FieldLabel>

                                    <Input
                                        id="username"
                                        name="username"
                                        type="text"
                                        defaultValue={auth?.attemptedUsername ?? ""}
                                        autoComplete="username"
                                        required
                                        className={cn(usernameError && "border-red-500")}
                                    />

                                    {usernameError && (
                                        <FieldDescription className="text-red-600">
                                            {messagesPerField.get("username")}
                                        </FieldDescription>
                                    )}
                                </Field>

                                {/* Password */}
                                <Field>
                                    <div className="flex items-center">
                                        <FieldLabel htmlFor="password">
                                            {i18n.msg("password")}
                                        </FieldLabel>

                                        {realm.resetPasswordAllowed && (
                                            <a
                                                href={url.loginResetCredentialsUrl}
                                                className="ml-auto inline-block text-sm underline-offset-4 hover:underline"
                                            >
                                                {i18n.msg("doForgotPassword")}
                                            </a>
                                        )}
                                    </div>

                                    <Input
                                        id="password"
                                        name="password"
                                        type="password"
                                        required
                                        autoComplete="current-password"
                                        className={cn(usernameError && "border-red-500")}
                                    />
                                </Field>

                                {/* Remember me */}
                                {realm.rememberMe && (
                                    <Field>
                                        <label htmlFor="rememberMe" className="flex items-center gap-x-2 cursor-pointer">
                                            <input
                                                type="checkbox"
                                                id="rememberMe"
                                                name="rememberMe"
                                                defaultChecked={realm.rememberMe}
                                                className="size-4"
                                            />
                                            <span className="text-sm">{i18n.msg("rememberMe")}</span>
                                        </label>
                                    </Field>
                                )}

                                {/* Submit */}
                                <Field>
                                    <Button type="submit" className="w-full">
                                        {i18n.msg("doLogIn")}
                                    </Button>
                                </Field>

                                {/* Social providers */}
                                {social?.providers?.map(provider => (
                                    <Field key={provider.providerId}>
                                        <Button
                                            variant="outline"
                                            type="submit"
                                            className="w-full"
                                            name="credentialId"
                                            value={provider.providerId}
                                        >
                                            {provider.displayName}
                                        </Button>
                                    </Field>
                                ))}

                                {/* Registration */}
                                {realm.registrationAllowed && (
                                    <Field>
                                        <FieldDescription className="text-center">
                                            {i18n.msg("noAccount")}{" "}
                                            <a
                                                href={url.registrationUrl}
                                                className="underline underline-offset-4"
                                            >
                                                {i18n.msg("doRegister")}
                                            </a>
                                        </FieldDescription>
                                    </Field>
                                )}
                            </FieldGroup>
                        </form>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}
