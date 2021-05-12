package manifest

func WithReleaseName(v string) ProcessFunc {
	return func(m Object) Object {
		// namespace shouldn't inject label
		if m.GetObjectKind().GroupVersionKind().Kind == "Namespace" {
			return m
		}

		labels := m.GetLabels()

		if labels == nil {
			labels = map[string]string{}
		}

		labels[LabelRelease] = v

		m.SetLabels(labels)

		return m
	}
}
