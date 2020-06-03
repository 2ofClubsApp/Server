import React from 'react'
import Button from "react-bootstrap/Button";
import {useHistory} from 'react-router-dom'


export const LandingPage = () => {
    const history = useHistory();
    const changeRoute = (path: string) => {
        history.replace({pathname: path})
    };

    return (
        <div>
            Welcome to 2ofClubs
            <Button onClick={() => changeRoute('/login')}>Login</Button>
        </div>
    );

};

