import React from 'react'
import Form from "react-bootstrap/Form";
import {FormLabel} from "../../types/FormInfo";
import {setUserData} from "../../store/actions/actions";
import {connect} from "react-redux"

let action = {}

const Label = (label: FormLabel) => {
    return (
        <Form.Group controlId={label.info.controlId}>
            <Form.Label>{label.info.label}</Form.Label>
            <Form.Control
                onChange={(event: React.ChangeEvent<HTMLTextAreaElement>) => action = setUserData(event.target.id, event.target.value)}
                type={label.info.type} placeholder={label.info.placeholder}/>
        </Form.Group>
    )
};

export default connect(null, action)(Label);
