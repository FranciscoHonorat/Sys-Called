{{- define "sys-called.labels" -}}
app.kubernetes.io/part-of: sys-called
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/instance: {{ .Release.Name }}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version }}
{{- end }}

{{- define "sys-called.selector" -}}
app.kubernetes.io/instance: {{ .root.Release.Name }}
app.kubernetes.io/name: {{ .name }}
{{- end }}

{{- define "sys-called.image" -}}
{{ printf "%s:%s" .repository (toString .tag) }}
{{- end }}

{{- define "sys-called.secretName" -}}
sys-called-secrets
{{- end }}

{{- define "sys-called.secretValue" -}}
{{- $existing := lookup "v1" "Secret" .root.Release.Namespace (include "sys-called.secretName" .root) -}}
{{- if .value -}}
{{ .value | b64enc }}
{{- else if and $existing (index $existing.data .key) -}}
{{ index $existing.data .key }}
{{- else -}}
{{ .generate | b64enc }}
{{- end -}}
{{- end }}

{{- define "sys-called.containerSecurityContext" -}}
allowPrivilegeEscalation: false
readOnlyRootFilesystem: true
runAsNonRoot: true
runAsUser: 10001
runAsGroup: 10001
capabilities:
  drop: ["ALL"]
{{- end }}

{{- define "sys-called.initContainers" -}}
- name: ensure-kafka-topic
  image: {{ include "sys-called.image" .root.Values.images.kafka }}
  imagePullPolicy: {{ .root.Values.imagePullPolicy }}
  command:
    - sh
    - -c
    - >-
      until /opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka:9092
      --create --if-not-exists --topic {{ .root.Values.kafka.topic }}
      --partitions 1 --replication-factor 1;
      do echo "waiting for kafka"; sleep 3; done
- name: copy-migrations
  image: {{ .image }}
  imagePullPolicy: {{ .root.Values.imagePullPolicy }}
  command: ["sh", "-c", "cp /migrations/*.sql /work/"]
  securityContext:
    {{- include "sys-called.containerSecurityContext" . | nindent 4 }}
  volumeMounts:
    - name: migrations
      mountPath: /work
- name: migrate
  image: {{ include "sys-called.image" .root.Values.images.postgres }}
  imagePullPolicy: {{ .root.Values.imagePullPolicy }}
  env:
    - name: PGHOST
      value: {{ .db }}-postgres
    - name: PGUSER
      value: {{ .database.user }}
    - name: PGDATABASE
      value: {{ .database.name }}
    - name: PGPASSWORD
      valueFrom:
        secretKeyRef:
          name: {{ include "sys-called.secretName" .root }}
          key: {{ .db }}-db-password
  command:
    - sh
    - -c
    - >-
      until pg_isready -q; do echo "waiting for postgres"; sleep 2; done;
      for f in /work/*.sql; do echo "applying $f";
      psql -v ON_ERROR_STOP=1 -q -f "$f" || exit 1; done
  volumeMounts:
    - name: migrations
      mountPath: /work
      readOnly: true
{{- end }}
