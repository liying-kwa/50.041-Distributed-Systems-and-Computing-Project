import React from 'react'

const SuccessCourse = props => (
    <tr>
      <td>{props.course.section} ({props.course.class})</td>
      <td>Success! This class has been added to your schedule</td>
      <td>Sucess</td>
    </tr>
)

export default SuccessCourse;