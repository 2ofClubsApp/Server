import React from 'react'
import Button from "react-bootstrap/Button";
import Label from "../components/Form/Label";
import Form from "../components/Form/Form";
import {userLabel, emailLabel, passLabel, passConfirmLabel} from "../types/FormInfo";
import {SignUpLabel} from "../types/User";
import {useHistory} from 'react-router-dom';
import '../app.css';
import axios from "../axios"

const Signup = () => {
    const history = useHistory();
    const changeRoute = (path: string) => {
        history.replace({pathname: path})
    };

    const labels = [userLabel, emailLabel, passLabel, passConfirmLabel];
    const formLabels = labels.map((label, index) => {
        return (
            <Label key={index} props={label}/>
        )
    });

    return (
        <div>
            <Button variant="outline-light" className="m-2 text-uppercase" onClick={() => changeRoute('/')}>Back to
                Home
            </Button>
            <Form formSubmit={signup} buttonName={SignUpLabel} labels={formLabels} title={SignUpLabel}/>
        </div>
    )
};

const signup = () => {
    axios.post("/signup").then(response => {
        console.log(response)
    }).catch(err => {
        console.log(err + "Unable to get student ;.;");
    })
    // console.log("Signing Up!")
};

export default Signup
