1. Add labels to cluster config so that it is known who ran the tests.
2. Update the config to the original one after any change during the tests.
3. Use a separate cluster for every test to avoid flake because of the order of test runs.
4. For a delete cluster test, ensure that cluster has been deleted from the cloud console as well, and for imported clusters, not.
5. Create a shared step to create Private GKE on Qase.
