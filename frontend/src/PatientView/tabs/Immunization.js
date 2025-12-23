import React, { useState } from 'react';
import './Immunization.css';

function Immunization() {
  const [immunization, setImmunization] = useState({
    // Children
    childrenNone: false,
    bcg: true,
    opv1: false,
    opv2: false,
    opv3: false,
    dpt1: true,
    dpt2: true,
    dpt3: true,
    measles: false,
    hepatitisB1: false,
    hepatitisB2: false,
    hepatitisB3: true,
    hepatitisA: true,
    varicellaChickenPox: true,
    // Young
    youngNone: false,
    hpv: false,
    mmr: true,
    // Pregnant
    pregnantNone: true,
    tetanusToxoid: false,
    // Elderly
    elderlyNone: true,
    pnuemococcalVaccine: false,
    fluVaccine: false,
    // Other
    otherImmunization: ''
  });

  return (
    <div className="immunization-content">
      <div className="immunization-sections">
        <div className="immunization-left">
          {/* Children Section */}
          <div className="immunization-section">
            <h4>1. Children</h4>
            <div className="immunization-grid">
              <label>
                <input
                  type="checkbox"
                  checked={immunization.childrenNone}
                  onChange={(e) => setImmunization({...immunization, childrenNone: e.target.checked})}
                />
                <span>None</span>
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={immunization.dpt1}
                  onChange={(e) => setImmunization({...immunization, dpt1: e.target.checked})}
                />
                <span>DPT1</span>
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={immunization.hepatitisB2}
                  onChange={(e) => setImmunization({...immunization, hepatitisB2: e.target.checked})}
                />
                <span>Hepatitis B2</span>
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={immunization.bcg}
                  onChange={(e) => setImmunization({...immunization, bcg: e.target.checked})}
                />
                <span>BCG</span>
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={immunization.dpt2}
                  onChange={(e) => setImmunization({...immunization, dpt2: e.target.checked})}
                />
                <span>DPT2</span>
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={immunization.hepatitisB3}
                  onChange={(e) => setImmunization({...immunization, hepatitisB3: e.target.checked})}
                />
                <span>Hepatitis B3</span>
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={immunization.opv1}
                  onChange={(e) => setImmunization({...immunization, opv1: e.target.checked})}
                />
                <span>OPV1</span>
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={immunization.dpt3}
                  onChange={(e) => setImmunization({...immunization, dpt3: e.target.checked})}
                />
                <span>DPT3</span>
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={immunization.hepatitisA}
                  onChange={(e) => setImmunization({...immunization, hepatitisA: e.target.checked})}
                />
                <span>Hepatitis A</span>
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={immunization.opv2}
                  onChange={(e) => setImmunization({...immunization, opv2: e.target.checked})}
                />
                <span>OPV2</span>
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={immunization.measles}
                  onChange={(e) => setImmunization({...immunization, measles: e.target.checked})}
                />
                <span>Measles</span>
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={immunization.varicellaChickenPox}
                  onChange={(e) => setImmunization({...immunization, varicellaChickenPox: e.target.checked})}
                />
                <span>Varicella (Chicken Pox)</span>
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={immunization.opv3}
                  onChange={(e) => setImmunization({...immunization, opv3: e.target.checked})}
                />
                <span>OPV3</span>
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={immunization.hepatitisB1}
                  onChange={(e) => setImmunization({...immunization, hepatitisB1: e.target.checked})}
                />
                <span>Hepatitis B1</span>
              </label>
            </div>
          </div>

          {/* Young Section */}
          <div className="immunization-section">
            <h4>2. Young</h4>
            <div className="immunization-grid">
              <label>
                <input
                  type="checkbox"
                  checked={immunization.youngNone}
                  onChange={(e) => setImmunization({...immunization, youngNone: e.target.checked})}
                />
                <span>None</span>
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={immunization.hpv}
                  onChange={(e) => setImmunization({...immunization, hpv: e.target.checked})}
                />
                <span>HPV</span>
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={immunization.mmr}
                  onChange={(e) => setImmunization({...immunization, mmr: e.target.checked})}
                />
                <span>MMR</span>
              </label>
            </div>
          </div>

          {/* Pregnant Section */}
          <div className="immunization-section">
            <h4>3. Pregnant</h4>
            <div className="immunization-grid">
              <label>
                <input
                  type="checkbox"
                  checked={immunization.pregnantNone}
                  onChange={(e) => setImmunization({...immunization, pregnantNone: e.target.checked})}
                />
                <span>None</span>
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={immunization.tetanusToxoid}
                  onChange={(e) => setImmunization({...immunization, tetanusToxoid: e.target.checked})}
                />
                <span>Tetanus Toxoid</span>
              </label>
            </div>
          </div>

          {/* Elderly Section */}
          <div className="immunization-section">
            <h4>4. Elderly</h4>
            <div className="immunization-grid">
              <label>
                <input
                  type="checkbox"
                  checked={immunization.elderlyNone}
                  onChange={(e) => setImmunization({...immunization, elderlyNone: e.target.checked})}
                />
                <span>None</span>
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={immunization.pnuemococcalVaccine}
                  onChange={(e) => setImmunization({...immunization, pnuemococcalVaccine: e.target.checked})}
                />
                <span>Pnuemococcal Vaccine</span>
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={immunization.fluVaccine}
                  onChange={(e) => setImmunization({...immunization, fluVaccine: e.target.checked})}
                />
                <span>Flu Vaccine</span>
              </label>
            </div>
          </div>
        </div>

        <div className="immunization-right">
          {/* Other Immunization Section */}
          <div className="immunization-section other-section">
            <h4>5. Other Immunization</h4>
            <textarea
              value={immunization.otherImmunization}
              onChange={(e) => setImmunization({...immunization, otherImmunization: e.target.value})}
              placeholder="lorem epsum 12345"
              rows="30"
            ></textarea>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Immunization;
