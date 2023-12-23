# Doka-connector
- **Apache Kafka에서 AWS Kinesis로 데이터 전송**

## 소개

이 프로젝트는 Apache Kafka의 메시지를 AWS Kinesis로 전송하는 커넥터를 구현합니다.
EKSAWS Elastic Kubernetes Service(EKS)를 사용하여 구축되었으며, 효율적인 데이터 스트리밍과 처리를 위해 AWS 서비스와 통합되어 있습니다.

## 주요 기능

- **MSK (Managed Streaming for Apache Kafka) 지원**: Kafka 특정 토픽의 메시지를 소비하고 처리합니다.
- **AWS Kinesis 통합**: 스트림 데이터를 AWS Kinesis로 전송합니다.
- **동적 자격 증명 갱신**: CrossAccount STS를 30분마다 자동으로 갱신합니다.

## 시작하기

### 전제 조건

- AWS 계정
- AWS CLI 설치 및 구성
- Kafka 클러스터 (MSK 또는 기타 Kafka 클러스터)
- AWS Kinesis Stream 설정

### 설치 방법

1. 리포지토리 Clone
2. 필요한 설정 파일 구성 (예: Kafka 브로커 정보, AWS 자격 증명, Kinesis 스트림 이름 등)
- 기본적으로 EKS에 배포하는것을 기본으로 만들었으며, 각 환경별로 MSK, Role ARN이 다르므로 {{ env }}.yaml에 적절한 값을 추가해야 합니다.
- go run main.go를 통해 Local에서 진행해도 되며, EKS에 배포해서 사용해도 됩니다.

## 구성 및 사용 방법
- `configs` 패키지: Kafka 및 AWS Kinesis 관련 설정을 관리합니다.
- `streamers` 패키지: AWS Kinesis와의 통신을 담당하며, STS를 통한 자격 증명 갱신 기능을 포함합니다.
- `consumers` 패키지: Kafka Consumer Group을 구성하고 메시지를 소비한 후 `streamers`를 통해 Kinesis에 전송합니다.
