import React from 'react'
import Button from "react-bootstrap/Button";
import Label from "../components/Form/Label";
import Form from "../components/Form/Form";
import {emailLabel, passLabel} from "../types/FormInfo";
import {useHistory} from 'react-router-dom'
import '../app.css';
import {LoginLabel} from "../types/User";


const Login = () => {
    const history = useHistory();
    const changeRoute = (path: string) => {
        history.replace({pathname: path})
    };

    const [state, setState] = React.useState({
        username: "",
        email: "",
        password: "",
    });

    const handleChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
        const value = event.target.value
        const id = event.target.id
        setState({
            ...state,
            [id]: value
        })
    }

    const labels = [emailLabel, passLabel];
    const formLabels = labels.map((label, index) => {
        return (
            <Label key={index} info={label}/>
        )
    });

    return (
        <div>
            <Button variant="outline-light" className="m-2 text-uppercase" onClick={() => changeRoute('/')}>Back to
                Home
            </Button>
            <Form formSubmit={changeRoute} labels={formLabels} title={LoginLabel} buttonName={LoginLabel}/>
        </div>
    )
};

// const login = () => {
//     console.log(emailLabel.controlId)
// };

export default Login
