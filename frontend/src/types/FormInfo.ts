export type FormInfo = {
    controlId: string
    label: string
    type: string
    placeholder: string
}
export const userLabel: FormInfo = {
    controlId: "username",
    label: "Username",
    type: "text",
    placeholder: "Username"
}
export const emailLabel: FormInfo = {
    controlId: "email",
    label: "Email Address",
    type: "email",
    placeholder: "Email"
};
export const passLabel: FormInfo = {
    controlId: "password",
    label: "Password",
    type: "password",
    placeholder: "Password"
};
export const passConfirmLabel: FormInfo = {
    controlId: "confirmpassword",
    label: "Confirm Password",
    type: "password",
    placeholder: "Password"
};
