import React from 'react';

// Sanitizers
const onlyDigits = (val) => (val || '').toString().replace(/\D+/g, '');
const onlyLetters = (val) => (val || '').toString().replace(/[^a-zA-Z\s]+/g, '');

const ItemEditForm = ({ item, onChange, onSave, onCancel, loading, sexOptions, civilStatusOptions }) => {
  return (
    <>
      <div className="form-row">
        <div className="form-field">
          <label>Case No.</label>
          <input
            name="caseNo"
            type="text"
            value={item.caseNo || ''}
            onChange={e => onChange({ ...item, caseNo: onlyDigits(e.target.value) })}
            placeholder="Case No."
            disabled={loading}
          />
        </div>
        <div className="form-field">
          <label>Hospital No.</label>
          <input
            name="hospitalNo"
            type="text"
            value={item.hospitalNo || ''}
            onChange={e => onChange({ ...item, hospitalNo: onlyDigits(e.target.value) })}
            placeholder="Hospital No."
            disabled={loading}
          />
        </div>
        <div className="form-field">
          <label>Age</label>
          <input
            name="age"
            type="text"
            value={item.age || ''}
            onChange={e => onChange({ ...item, age: onlyDigits(e.target.value) })}
            placeholder="Age"
            disabled={loading}
          />
        </div>
        <div className="form-field">
          <label>Room</label>
          <input
            name="room"
            type="text"
            value={item.room || ''}
            onChange={e => onChange({ ...item, room: onlyDigits(e.target.value) })}
            placeholder="Room"
            disabled={loading}
          />
        </div>
      </div>
      <div className="form-row">
        <div className="form-field">
          <label>Last Name</label>
          <input
            name="lastname"
            type="text"
            value={item.lastname}
            onChange={e => onChange({ ...item, lastname: onlyLetters(e.target.value) })}
            placeholder="Lastname"
            disabled={loading}
            required
          />
        </div>
        <div className="form-field">
          <label>First Name</label>
          <input
            name="firstname"
            type="text"
            value={item.firstname}
            onChange={e => onChange({ ...item, firstname: onlyLetters(e.target.value) })}
            placeholder="Firstname"
            disabled={loading}
            required
          />
        </div>
        <div className="form-field">
          <label>Middle Name</label>
          <input
            name="middlename"
            type="text"
            value={item.middlename}
            onChange={e => onChange({ ...item, middlename: onlyLetters(e.target.value) })}
            placeholder="Middlename"
            disabled={loading}
          />
        </div>
        <div className="form-field">
          <label>Suffix</label>
          <input
            name="suffix"
            type="text"
            value={item.suffix}
            onChange={e => onChange({ ...item, suffix: onlyLetters(e.target.value) })}
            placeholder="Suffix"
            disabled={loading}
          />
        </div>
      </div>
      <div className="form-row">
        <div className="form-field">
          <label>Birthdate</label>
          <input
            name="birthdate"
            type="date"
            value={item.birthdate}
            onChange={e => onChange({ ...item, birthdate: e.target.value })}
            placeholder="Birthdate"
            disabled={loading}
          />
        </div>
        <div className="form-field">
          <label>Sex</label>
          <select
            name="sex"
            value={item.sex}
            onChange={e => onChange({ ...item, sex: e.target.value })}
            disabled={loading}
            required
          >
            {sexOptions.map(opt => <option key={opt} value={opt}>{opt || 'Sex'}</option>)}
          </select>
        </div>
        <div className="form-field">
          <label>Civil Status</label>
          <select
            name="civil_status"
            value={item.civil_status}
            onChange={e => onChange({ ...item, civil_status: e.target.value })}
            disabled={loading}
            required
          >
            {civilStatusOptions.map(opt => <option key={opt} value={opt}>{opt || 'Civil Status'}</option>)}
          </select>
        </div>
      </div>
      <div className="form-row">
        <div className="form-field">
          <label>Admission Date</label>
          <input
            name="admissionDate"
            type="datetime-local"
            value={item.admissionDate || ''}
            onChange={e => onChange({ ...item, admissionDate: e.target.value })}
            placeholder="Admission Date"
            disabled={loading}
          />
        </div>
        <div className="form-field">
          <label>Discharge Date</label>
          <input
            name="dischargeDate"
            type="datetime-local"
            value={item.dischargeDate || ''}
            onChange={e => onChange({ ...item, dischargeDate: e.target.value })}
            placeholder="Discharge Date"
            disabled={loading}
          />
        </div>
        <div className="form-field">
          <label>Height (cm)</label>
          <input
            name="height"
            type="text"
            value={item.height || ''}
            onChange={e => onChange({ ...item, height: onlyDigits(e.target.value) })}
            placeholder="Height (cm)"
            disabled={loading}
          />
        </div>
        <div className="form-field">
          <label>Weight (kg)</label>
          <input
            name="weight"
            type="text"
            value={item.weight || ''}
            onChange={e => onChange({ ...item, weight: onlyDigits(e.target.value) })}
            placeholder="Weight (kg)"
            disabled={loading}
          />
        </div>
      </div>
      <div className="form-row">
        <div className="form-field" style={{ width: '100%' }}>
          <label>Complaint</label>
          <textarea
            name="complaint"
            value={item.complaint || ''}
            onChange={e => onChange({ ...item, complaint: onlyLetters(e.target.value) })}
            placeholder="Complaint"
            disabled={loading}
            rows={3}
            style={{ width: '100%' }}
          />
        </div>
      </div>
      <div className="button-group">
        <button onClick={onSave} disabled={loading} className="save-btn">Save</button>
        <button onClick={onCancel} disabled={loading} className="cancel-btn">Cancel</button>
      </div>
    </>
  );
};

export default ItemEditForm;
