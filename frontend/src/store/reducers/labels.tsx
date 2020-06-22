import {actionDefinition, SET_USERDATA} from "../actions/actionTypes";
import {User} from "../types/users/types";

const initialState: User = {
    username: "",
    email: "",
    password: "",
    confirmPassword: ""
}

export default function (state = initialState, action: actionDefinition) {
    console.log(action)
    switch (action.type) {
        case SET_USERDATA:
            console.log("bananas");
            break;
        default:
            return state;
    }

}
