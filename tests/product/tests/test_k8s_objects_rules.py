"""
Kubernetes CIS rules verification.
This module verifies correctness of retrieved findings by manipulating audit and remediation actions
"""
from datetime import datetime

import pytest

# from product.tests.data.k8s_object.k8s_object_rules import *
from product.tests.data.k8s_object import k8s_object_rules as k8s_tc
from commonlib.utils import get_evaluation, get_resource_identifier
from commonlib.framework.reporting import skip_param_case, SkipReportData


@pytest.mark.k8s_object_rules
@pytest.mark.parametrize(
    ("rule_tag", "resource_type", "resource_body", "expected"),
    [
        *k8s_tc.cis_5_1_3.values(),
        *k8s_tc.cis_5_1_5.values(),
        *k8s_tc.cis_5_1_6.values(),
        *k8s_tc.cis_5_2_3.values(),
        *k8s_tc.cis_5_2_4.values(),
        *k8s_tc.cis_5_2_5.values(),
        *skip_param_case(skip_list=[*k8s_tc.cis_5_1_5_skip.values(),
                                    *k8s_tc.cis_5_2_2.values(),
                                    *k8s_tc.cis_5_2_3_skip.values(),
                                    *k8s_tc.cis_5_2_4_skip.values(),
                                    *k8s_tc.cis_5_2_5_skip.values(),
                                    *k8s_tc.cis_5_2_6.values(),
                                    *k8s_tc.cis_5_2_7.values(),
                                    *k8s_tc.cis_5_2_8.values()
                                    ],
                         data_to_report=SkipReportData(
                             url_title="security-team: #4312",
                             url_link="https://github.com/elastic/security-team/issues/4312",
                             skip_reason="known issue: broken k8s object tests"
                         ))
        # *k8s_tc.cis_5_2_9.values(), - TODO: cases are not implemented
        # *k8s_tc.cis_5_2_10.values() - TODO: cases are not implemented
    ],
    ids=[
        *k8s_tc.cis_5_1_3.keys(),
        *k8s_tc.cis_5_1_5.keys(),
        *k8s_tc.cis_5_1_6.keys(),
        *k8s_tc.cis_5_2_3.keys(),
        *k8s_tc.cis_5_2_4.keys(),
        *k8s_tc.cis_5_2_5.keys(),
        *k8s_tc.cis_5_1_5_skip.keys(),
        *k8s_tc.cis_5_2_2.keys(),
        *k8s_tc.cis_5_2_3_skip.keys(),
        *k8s_tc.cis_5_2_4_skip.keys(),
        *k8s_tc.cis_5_2_5_skip.keys(),
        *k8s_tc.cis_5_2_6.keys(),
        *k8s_tc.cis_5_2_7.keys(),
        *k8s_tc.cis_5_2_8.keys(),
        # *k8s_tc.cis_5_2_9.keys(), - TODO: cases are not implemented
        # *k8s_tc.cis_5_2_10.keys() - TODO: cases are not implemented
    ]
)
def test_kube_resource_patch(test_env, rule_tag, resource_type, resource_body, expected):
    """
    Test kube resource
    @param test_env: pre step that set-ups a kube resources to test on
    @param rule_tag: rule tag in the CIS benchmark
    @param resource_type: kube resource type, e.g., Pod, ServiceAccount, etc.
    @param resource_body: a dict to represent the relevant properties of the resource
    @param expected: "failed" or "passed"
    """
    k8s_client, _, agent_config = test_env

    # make sure resource exists
    metadata = resource_body['metadata']
    relevant_metadata = {k: metadata[k] for k in ('name', 'namespace') if k in metadata}
    try:
        resource = k8s_client.get_resource(resource_type=resource_type, **relevant_metadata)
    except TypeError as type_error:
        print(type_error)
        resource = k8s_client.get_resource(resource_type=resource_type,
                                           namespace=agent_config.namespace,
                                           **relevant_metadata)

    assert resource, f"Resource {resource_type} not found"

    # patch resource
    resource = k8s_client.patch_resources(
        resource_type=resource_type,
        body=resource_body,
        **relevant_metadata
    )
    if resource is None:
        raise ValueError(
            f'Could not patch resource type {resource_type}:'
            f' {relevant_metadata} with patch {resource_body}')

    # check resource evaluation
    pods = k8s_client.get_agent_pod_instances(agent_name=agent_config.name,
                                              namespace=agent_config.namespace)

    evaluation = get_evaluation(
        k8s=k8s_client,
        timeout=agent_config.findings_timeout,
        pod_name=pods[0].metadata.name,
        namespace=agent_config.namespace,
        rule_tag=rule_tag,
        exec_timestamp=datetime.utcnow(),
        resource_identifier=get_resource_identifier(resource_body),
    )

    assert evaluation is not None, f"No evaluation for rule {rule_tag} could be found"
    assert evaluation == expected, f"Rule {rule_tag} verification failed, " \
                                   f"expected: {expected} actual: {evaluation}"
