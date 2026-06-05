# fhir-test-discovery-example

```bash
yourphr git:(docs/fhir-testing) ✗ curl -s https://fhir.fhirpoint.open.allscripts.com/fhirroute/open/76308/.well-known/smart-configuration | python3 -m json.tool
```

Returns:

```json
{
    "authorization_endpoint": "https://open.allscripts.com/fhirroute/fmhpatientauth/fmhorgid/b065c577-9909-41ab-9c77-a40600f66d54/connect/authorize",
    "issuer": "https://muauthentication.followmyhealth.com",
    "jwks_uri": "https://muauthentication.followmyhealth.com/api/jwks",
    "token_endpoint": "https://muauthentication.followmyhealth.com/api/access",
    "token_endpoint_auth_methods": [
        "client_secret_post",
        "client_secret_basic"
    ],
    "scopes_supported": [
        "openid",
        "profile",
        "launch/patient",
        "offline_access",
        "fhiruser",
        "patient/Medication.read",
        "patient/MedicationOrder.read",
        "patient/MedicationStatement.read",
        "patient/MedicationRequest.read",
        "patient/MedicationDispense.read",
        "patient/Practitioner.read",
        "patient/PractitionerRole.read",
        "patient/DiagnosticReport.read",
        "patient/DiagnosticOrder.read",
        "patient/DocumentReference.read",
        "patient/Binary.read",
        "patient/Composition.read",
        "patient/AllergyIntolerance.read",
        "patient/CarePlan.read",
        "patient/CareTeam.read",
        "patient/Condition.read",
        "patient/Coverage.read",
        "patient/Patient.read",
        "patient/Encounter.read",
        "patient/Goal.read",
        "patient/Immunization.read",
        "patient/Device.read",
        "patient/Location.read",
        "patient/Observation.read",
        "patient/Organization.read",
        "patient/Procedure.read",
        "patient/Provenance.read",
        "patient/QuestionnaireResponse.read",
        "patient/RelatedPerson.read",
        "patient/ServiceRequest.read",
        "patient/Specimen.read",
        "patient/DocumentReference.write"
    ],
    "response_types_supported": [
        "code",
        "token",
        "code id_token",
        "code token",
        "code id_token token"
    ],
    "grant_types_supported": [
        "authorization_code",
        "implicit",
        "hybrid"
    ],
    "code_challenge_methods_supported": [
        "S256"
    ],
    "management_endpoint": null,
    "capabilities": [
        "launch-ehr",
        "launch-standalone",
        "client-public",
        "client-confidential-symmetric",
        "context-ehr-patient",
        "context-ehr-encounter",
        "context-standalone-patient",
        "context-standalone-encounter",
        "context-passthrough-banner",
        "context-passthrough-style",
        "context-banner",
        "context-style",
        "sso-openid-connect",
        "permission-offline",
        "permission-patient",
        "permission-user",
        "client-confidential-asymmetric",
        "authorize-post",
        "permission-v1",
        "permission-v2"
    ]
}
```
