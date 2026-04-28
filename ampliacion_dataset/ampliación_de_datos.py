import os
import random
import ipaddress
from datetime import timedelta

import pandas as pd


BASE_DIR = os.path.dirname(os.path.abspath(__file__))
INPUT_FILE = os.path.join(BASE_DIR, "../dataset/dataset_base.csv")
OUTPUT_FILE = os.path.join(BASE_DIR, "../dataset/dataset_1M_raw.csv")

TARGET_SIZE = 1_000_000
MIN_TEXTO_LEN = 28

# Probabilidades
CLUSTER_SPAM_PERCENT = 0.02
SPAM_PERCENT = 0.02
NOISE_PERCENT = 0.08
DUPLICATE_PERCENT = 0.02

INVALID_IP_PERCENT = 0.02
INVALID_TIMESTAMP_PERCENT = 0.02
EMPTY_FIELD_PERCENT = 0.01

NUM_USUARIOS_POOL = 300000

random.seed(42)


TEXTOS_BASE = [
    "No respetaron la garantia del producto y no ofrecieron cambio.",
    "El servicio fue deficiente y no resolvieron mi reclamo a tiempo.",
    "No responden a mis llamadas ni mensajes desde hace varios dias.",
    "Me cobraron de mas en el recibo y no explican el motivo.",
    "El producto llego defectuoso y no quieren aplicar la garantia.",
    "No cumplen con lo ofrecido en la publicidad del servicio.",
]

SUFIJOS_CONTEXTO = [
    "Necesito una solucion formal en el menor tiempo posible.",
    "Solicito atencion inmediata porque ya presente reclamo previo.",
    "Adjuntare evidencia del caso para validar lo reportado.",
    "Espero respuesta por canal oficial y devolucion correspondiente.",
]


def generar_usuario():
    return f"USR{random.randint(1, NUM_USUARIOS_POOL):06d}"


def generar_ip():
    return str(ipaddress.IPv4Address(random.randint(0x0B000000, 0xDF000000)))


def generar_ip_invalida():
    return random.choice([
        "999.999.999.999",
        "abc.def.ghi.jkl",
        "256.256.256.256",
        "",
        "IP_INVALIDA",
    ])


def generar_timestamp_invalido():
    return random.choice([
        "32/13/2016 99:99",
        "fecha_invalida",
        "",
        "2015-99-99 25:61:61",
    ])


def variar_timestamp(ts):
    if pd.isna(ts):
        return ts
    return ts + timedelta(
        days=random.randint(-3, 3),
        hours=random.randint(-12, 12),
        minutes=random.randint(-59, 59),
    )


def garantizar_texto_informativo(texto):
    texto = " ".join(str(texto).strip().split())

    if texto == "" or texto == "-":
        texto = random.choice(TEXTOS_BASE)

    if len(texto) < MIN_TEXTO_LEN:
        texto = f"{texto} {random.choice(SUFIJOS_CONTEXTO)}"

    return " ".join(texto.split())


def mutar_texto(texto):
    texto = garantizar_texto_informativo(texto)

    variantes = [
        texto,
        texto.lower(),
        texto.upper(),
        f"{texto}!!!",
        f"IMPORTANTE: {texto}",
        f"{texto} POR FAVOR",
        texto.replace("a", "@", 1),
    ]

    candidato = random.choice(variantes)
    return garantizar_texto_informativo(candidato)


def generar_cluster_spam(base, size=20):
    cluster = []

    texto_spam = random.choice([
        "RECLAMO NO RESUELTO HACE SEMANAS, NECESITO SOLUCION INMEDIATA DEL SERVICIO.",
        "NO RESPONDEN A MIS LLAMADAS NI CORREOS, EXIJO UNA RESPUESTA FORMAL AHORA.",
        "ESTAFA EN LA FACTURACION, COBRAN MONTOS INDEBIDOS Y NO DAN EXPLICACION.",
        "URGE DEVOLUCION DEL DINERO POR INCUMPLIMIENTO DEL SERVICIO OFRECIDO.",
        "NO ATENDIERON MI GARANTIA Y EL PRODUCTO SIGUE FALLANDO DESDE LA COMPRA.",
    ])

    ip_fija = generar_ip()
    usuario_fijo = generar_usuario()
    base_time = base["timestamp"]

    for _ in range(size):
        nuevo = base.copy()
        nuevo["texto_reclamo"] = mutar_texto(texto_spam)
        nuevo["ip_address"] = ip_fija
        nuevo["usuario_id"] = usuario_fijo

        if not pd.isna(base_time):
            nuevo["timestamp"] = base_time + timedelta(seconds=random.randint(0, 60))
        else:
            nuevo["timestamp"] = base_time

        nuevo["is_synthetic_spam"] = 1
        cluster.append(nuevo)

    return cluster


# =========================
# CARGA DATA
# =========================
df = pd.read_csv(INPUT_FILE, sep=";", encoding="utf-8", engine="python")
df = df.loc[:, df.columns.str.strip() != ""]
df.columns = df.columns.str.strip()

print("Columnas detectadas:", df.columns.tolist())

if "timestamp" not in df.columns:
    df["timestamp"] = pd.to_datetime(df["FECHA_PRESENTACION"], errors="coerce", dayfirst=True)
else:
    df["timestamp"] = pd.to_datetime(df["timestamp"], errors="coerce")

if "texto_reclamo" not in df.columns:
    df["texto_reclamo"] = [random.choice(TEXTOS_BASE) for _ in range(len(df))]

if "usuario_id" not in df.columns:
    df["usuario_id"] = "USR000000"

if "ip_address" not in df.columns:
    df["ip_address"] = "0.0.0.0"

df["texto_reclamo"] = df["texto_reclamo"].apply(garantizar_texto_informativo)
base_records = df.to_dict("records")
registros = base_records.copy()

next_id = len(registros) + 1

while len(registros) < TARGET_SIZE:
    base = random.choice(base_records)

    if random.random() < CLUSTER_SPAM_PERCENT:
        cluster = generar_cluster_spam(base, size=random.randint(10, 40))
        for nuevo in cluster:
            nuevo["id_reclamo"] = next_id
            next_id += 1
            registros.append(nuevo)
        continue

    nuevo = base.copy()
    nuevo["usuario_id"] = generar_usuario()
    nuevo["ip_address"] = generar_ip()
    nuevo["timestamp"] = variar_timestamp(base["timestamp"])

    if random.random() < DUPLICATE_PERCENT:
        nuevo["id_reclamo"] = next_id
        next_id += 1
        registros.append(nuevo)
        continue

    if random.random() < SPAM_PERCENT:
        nuevo["texto_reclamo"] = mutar_texto(random.choice(TEXTOS_BASE))
        nuevo["is_synthetic_spam"] = 1
    else:
        nuevo["texto_reclamo"] = mutar_texto(base.get("texto_reclamo", random.choice(TEXTOS_BASE)))
        nuevo["is_synthetic_spam"] = 0

    if random.random() < NOISE_PERCENT:
        nuevo["texto_reclamo"] = mutar_texto(nuevo["texto_reclamo"])

    if random.random() < INVALID_IP_PERCENT:
        nuevo["ip_address"] = generar_ip_invalida()

    if random.random() < INVALID_TIMESTAMP_PERCENT:
        nuevo["timestamp"] = generar_timestamp_invalido()

    if random.random() < EMPTY_FIELD_PERCENT:
        # Evitamos vaciar texto para no introducir advertencias artificiales de corto/vacio.
        nuevo["usuario_id"] = ""

    nuevo["id_reclamo"] = next_id
    next_id += 1
    registros.append(nuevo)

    if len(registros) % 50000 == 0:
        print(f"Generados: {len(registros):,}")


df_final = pd.DataFrame(registros)
df_final["timestamp"] = df_final["timestamp"].astype(str)
df_final.to_csv(OUTPUT_FILE, index=False)

print(f"\nDataset generado: {len(df_final):,} registros")
