import os
import random
from datetime import timedelta

import pandas as pd

# =========================================================
# CONFIG
# =========================================================

BASE_DIR = os.path.dirname(os.path.abspath(__file__))

INPUT_FILE = os.path.join(
    BASE_DIR,
    "../dataset/dataset_base.csv",
)

OUTPUT_FILE = os.path.join(
    BASE_DIR,
    "../dataset/dataset_1M_raw.csv",
)

LOGS_DIR = os.path.join(
    BASE_DIR,
    "logs",
)

SUMMARY_FILE = os.path.join(
    LOGS_DIR,
    "generation_summary.txt",
)

TARGET_SIZE = 1_700_000
CHUNK_SIZE = 50_000

NUM_USUARIOS_POOL = 300_000

MIN_TEXTO_LEN = 30
MAX_TEXTO_LEN = 250

# =========================================================
# PROBABILIDADES
# =========================================================

SPAM_PERCENT = 0.02
CLUSTER_SPAM_PERCENT = 0.003

DIRTY_PERCENT = 0.30
NOISE_PERCENT = 0.08

DUPLICATE_PERCENT = 0.02
VARIATION_PERCENT = 0.20

INVALID_IP_PERCENT = 0.02
INVALID_TIMESTAMP_PERCENT = 0.02
EMPTY_FIELD_PERCENT = 0.01

random.seed(42)

# aliases rápidos
rc = random.choice
rr = random.random
ri = random.randint

# =========================================================
# CONSTANTES
# =========================================================

STRIP_CHARS = ".,!?:;"

CATEGORIAS = (
    "facturacion",
    "delivery",
    "garantia",
    "atencion_cliente",
    "servicios",
    "banca",
)

# =========================================================
# CONECTORES
# =========================================================

CONECTORES = (
    "porque",
    "además",
    "sin embargo",
    "por eso",
    "aunque",
    "pero",
    "entonces",
    "debido a que",
    "por lo tanto",
)

# =========================================================
# RECLAMOS LEGÍTIMOS
# =========================================================

RECLAMOS_LEGITIMOS = {
    "facturacion": {
        "problemas": [
            "Me cobraron",
            "Hay un cargo indebido",
            "Aparece un cobro duplicado",
            "Me descontaron más de lo debido",
            "El monto no coincide con lo acordado",
        ],
        "detalles": [
            "en el recibo de este mes",
            "en mi factura",
            "en la última transacción",
            "en la cuota del servicio",
            "en mi cuenta",
        ],
        "contexto": [
            "No encuentro justificación para este cargo",
            "Contraté por otro monto",
            "Revisé mi historial",
            "Esto no corresponde",
            "No autorizo este cobro",
        ],
        "resolucion": [
            "Necesito que revisen el cargo",
            "Solicito una aclaración",
            "Quiero devolución del dinero",
            "Requiero explicación",
            "Pido corregir la facturación",
        ],
    },

    "delivery": {
        "problemas": [
            "El paquete nunca llegó",
            "Mi pedido llegó tarde",
            "El producto llegó dañado",
            "El envío fue rechazado",
            "El producto no coincide",
        ],
        "detalles": [
            "hace más de una semana",
            "según el plazo acordado",
            "con golpes en la caja",
            "sin respuesta de la empresa",
            "con información incorrecta",
        ],
        "contexto": [
            "Necesitaba el producto en esa fecha",
            "Pagué por envío garantizado",
            "No es la primera vez",
            "El tracking no se actualiza",
            "Ya envié varias reclamaciones",
        ],
        "resolucion": [
            "Solicito reembolso completo",
            "Necesito un cambio urgente",
            "Exijo seguimiento del reclamo",
            "Pido compensación",
            "Quiero una solución inmediata",
        ],
    },

    "garantia": {
        "problemas": [
            "El producto falla constantemente",
            "El defecto es de fábrica",
            "No respetan la cobertura",
            "Rechazaron mi garantía",
            "El producto sigue roto",
        ],
        "detalles": [
            "aunque está dentro del plazo",
            "según el certificado",
            "sin explicación válida",
            "a pesar del comprobante",
            "más de una vez",
        ],
        "contexto": [
            "Compré hace poco",
            "Tengo el recibo original",
            "Esto ocurrió desde la compra",
            "No fue culpa mía",
            "Otros clientes reportan lo mismo",
        ],
        "resolucion": [
            "Exijo reparación inmediata",
            "Necesito garantía válida",
            "Solicito devolución del dinero",
            "Pido compensación",
            "Quiero un producto nuevo",
        ],
    },

    "atencion_cliente": {
        "problemas": [
            "No me atienden",
            "Llevo semanas esperando",
            "El personal fue grosero",
            "No resuelven mi problema",
            "Me transfieren constantemente",
        ],
        "detalles": [
            "cada vez que llamo",
            "desde hace días",
            "sin justificación",
            "a pesar de mis intentos",
            "en múltiples ocasiones",
        ],
        "contexto": [
            "Necesito ayuda urgente",
            "He llamado varias veces",
            "Otros clientes tienen el mismo problema",
            "El chat tampoco responde",
            "La atención es deficiente",
        ],
        "resolucion": [
            "Solicito atención inmediata",
            "Necesito hablar con un supervisor",
            "Exijo respuesta formal",
            "Requiero seguimiento",
            "Pido solución urgente",
        ],
    },

    "servicios": {
        "problemas": [
            "La velocidad de internet es muy baja",
            "La conexión cae constantemente",
            "El servicio tiene mala calidad",
            "Pagué por mejor servicio",
            "La cobertura es deficiente",
        ],
        "detalles": [
            "mucho menor a lo contratado",
            "todos los días",
            "según la publicidad",
            "sin solución",
            "en mi zona",
        ],
        "contexto": [
            "Afecta mis actividades",
            "Pago por servicio premium",
            "Esto lleva semanas",
            "No es la primera queja",
            "Necesito estabilidad",
        ],
        "resolucion": [
            "Solicito devolución",
            "Necesito mejorar el servicio",
            "Quiero descuento",
            "Exijo solución",
            "Pido cambio de plan",
        ],
    },

    "banca": {
        "problemas": [
            "Hubo un cargo no autorizado",
            "Mi transacción desapareció",
            "La transferencia nunca llegó",
            "Me cobraron comisiones indebidas",
            "Hay movimientos sospechosos",
        ],
        "detalles": [
            "en mi cuenta bancaria",
            "sin explicación",
            "en el estado de cuenta",
            "de empresas desconocidas",
            "no registrados por mí",
        ],
        "contexto": [
            "No autorizo esos movimientos",
            "Esto afecta mi seguridad",
            "Necesito proteger mi cuenta",
            "Temo fraude",
            "Es prioritario resolver esto",
        ],
        "resolucion": [
            "Solicito devolución inmediata",
            "Pido bloqueo de la cuenta",
            "Exijo investigación",
            "Necesito revisión urgente",
            "Requiero cambio de contraseña",
        ],
    },
}

# =========================================================
# SPAM
# =========================================================

TIPOS_SPAM = (
    "urgencia",
    "repetitivo",
    "flood",
    "promocional",
    "semi_humano",
)

SPAM_FRASES_URGENCIA = (
    "URGENTE",
    "INMEDIATO",
    "YA",
    "AHORA",
    "PRIORITARIO",
)

SPAM_FRAGMENTOS = (
    "quiero respuesta ya",
    "solucion urgente",
    "nadie responde",
    "esto es una estafa",
    "respuesta inmediata",
    "haz click ahora",
    "gana dinero rapido",
    "premio disponible",
)

# =========================================================
# RUIDO
# =========================================================

PALABRAS_RUIDO = (
    "producto",
    "servicio",
    "cobro",
    "pago",
    "llamada",
    "respuesta",
    "reclamo",
    "devolucion",
    "dinero",
    "garantia",
    "factura",
    "problema",
    "error",
)

# =========================================================
# HELPERS
# =========================================================

def generar_usuario():
    return f"USR{ri(1, NUM_USUARIOS_POOL):06d}"


def generar_ip():
    return ".".join(
        str(ri(1, 254))
        for _ in range(4)
    )


def generar_ip_invalida():
    return rc((
        "999.999.999.999",
        "abc.def.ghi.jkl",
        "256.256.256.256",
        "",
        "IP_INVALIDA",
    ))


def generar_timestamp_invalido():
    return rc((
        "32/13/2016 99:99",
        "fecha_invalida",
        "",
        "2015-99-99 25:61:61",
    ))


def variar_timestamp(ts):

    if pd.isna(ts):
        return ts

    return ts + timedelta(
        days=ri(-3, 3),
        hours=ri(-12, 12),
        minutes=ri(-59, 59),
    )


# =========================================================
# TEXTO
# =========================================================

def insertar_conector(oracion1, oracion2):

    return (
        f"{oracion1} "
        f"{rc(CONECTORES)} "
        f"{oracion2}"
    )


def aplicar_variacion_casos(texto):

    caso = random.choices(
        ["upper", "lower", "mixed"],
        weights=[30, 30, 40],
        k=1
    )[0]

    if caso == "upper":
        return texto.upper()

    if caso == "lower":
        return texto.lower()

    oraciones = texto.split(". ")

    resultado = []

    for o in oraciones:

        if not o.strip():
            continue

        primera = o[0].upper()
        resto = o[1:]

        if rr() < 0.3:
            resto = resto.lower()

        resultado.append(
            primera + resto
        )

    return ". ".join(resultado)


# =========================================================
# RECLAMO LEGITIMO
# =========================================================

def generar_reclamo_legitimo():

    categoria = rc(CATEGORIAS)

    plantilla = RECLAMOS_LEGITIMOS[categoria]

    p1 = (
        f"{rc(plantilla['problemas'])} "
        f"{rc(plantilla['detalles'])}"
    )

    p2 = rc(plantilla["contexto"])

    p3 = rc(plantilla["resolucion"])

    texto = (
        insertar_conector(p1, p2)
        + ". "
        + insertar_conector(p2, p3)
    )

    return aplicar_variacion_casos(texto)


# =========================================================
# TEXTO INCOHERENTE
# =========================================================

def generar_texto_incoherente():

    tipo = ri(1, 4)

    if tipo == 1:

        return " ".join(
            rc(PALABRAS_RUIDO)
            for _ in range(ri(4, 8))
        )

    if tipo == 2:

        return " ".join(
            rc(PALABRAS_RUIDO)
            for _ in range(ri(8, 15))
        )

    if tipo == 3:

        return (
            f"{rc(PALABRAS_RUIDO)} "
            f"... "
            f"{rc(PALABRAS_RUIDO)}"
        )

    return (
        f"{rc(PALABRAS_RUIDO)} "
        f"{rc(PALABRAS_RUIDO)}"
    )


# =========================================================
# TEXTO GENERAL
# =========================================================

def generar_texto():

    if rr() < DIRTY_PERCENT:
        return generar_texto_incoherente()

    return generar_reclamo_legitimo()


# =========================================================
# VARIACION LEGITIMA
# =========================================================

def generar_variacion_legitima(texto):

    extras = (
        "Por favor revisen esto.",
        "Necesito una solución urgente.",
        "Tengo comprobantes disponibles.",
        "Esto me está causando problemas.",
    )

    extra = rc(extras)

    if len(texto) + len(extra) < MAX_TEXTO_LEN:
        texto = f"{texto} {extra}"

    return aplicar_variacion_casos(texto)


# =========================================================
# FEATURES NLP
# =========================================================

def calcular_caps_ratio(texto):

    letras = [
        c
        for c in texto
        if c.isalpha()
    ]

    if not letras:
        return 0.0

    mayus = sum(
        1
        for c in letras
        if c.isupper()
    )

    return mayus / len(letras)


def calcular_coherencia_basica(texto):

    palabras = texto.lower().split()

    if len(palabras) < 5:
        return False

    tiene_conector = any(
        c in texto.lower()
        for c in CONECTORES
    )

    return tiene_conector


# =========================================================
# SPAM
# =========================================================

def generar_spam_robotizado():

    tipo = rc(TIPOS_SPAM)

    if tipo == "urgencia":

        texto = (
            f"{rc(SPAM_FRASES_URGENCIA)} "
            f"{rc(SPAM_FRAGMENTOS)} "
            f"{rc(SPAM_FRASES_URGENCIA)}"
        )

    elif tipo == "repetitivo":

        palabra = rc((
            "URGENTE",
            "ESTAFA",
            "SOLUCION",
        ))

        texto = " ".join(
            palabra
            for _ in range(ri(5, 10))
        )

    elif tipo == "flood":

        texto = ". ".join(
            rc(SPAM_FRAGMENTOS)
            for _ in range(ri(5, 10))
        )

    elif tipo == "promocional":

        texto = (
            "GANA DINERO AHORA. "
            "CLICK AQUI. "
            "PREMIO DISPONIBLE."
        )

    else:

        texto = (
            f"{rc(SPAM_FRAGMENTOS)} "
            f"{rc(CONECTORES)} "
            f"{rc(SPAM_FRAGMENTOS)}"
        )

    return aplicar_variacion_casos(texto), tipo


# =========================================================
# REGISTRO
# =========================================================

def generar_registro(base_row, record_id):

    registro = dict(base_row)

    registro["id_reclamo"] = record_id
    registro["usuario_id"] = generar_usuario()
    registro["ip_address"] = generar_ip()

    registro["timestamp"] = variar_timestamp(
        registro["timestamp"]
    )

    registro["tipo_duplicado"] = "nuevo"

    # =====================================
    # SPAM
    # =====================================

    if rr() < SPAM_PERCENT:

        texto_spam, tipo_spam = (
            generar_spam_robotizado()
        )

        registro["texto_reclamo"] = texto_spam
        registro["spam_tipo"] = tipo_spam
        registro["is_synthetic_spam"] = 1

    else:

        texto = generar_texto()

        if rr() < NOISE_PERCENT:
            texto = generar_variacion_legitima(texto)

        registro["texto_reclamo"] = texto
        registro["spam_tipo"] = None
        registro["is_synthetic_spam"] = 0

    # =====================================
    # DUPLICADOS
    # =====================================

    if rr() < DUPLICATE_PERCENT:

        registro["tipo_duplicado"] = (
            "duplicado_exacto"
        )

    elif rr() < VARIATION_PERCENT:

        registro["tipo_duplicado"] = (
            "variacion"
        )

    # =====================================
    # INVALID DATA
    # =====================================

    if rr() < INVALID_IP_PERCENT:
        registro["ip_address"] = (
            generar_ip_invalida()
        )

    if rr() < INVALID_TIMESTAMP_PERCENT:
        registro["timestamp"] = (
            generar_timestamp_invalido()
        )

    if rr() < EMPTY_FIELD_PERCENT:
        registro["usuario_id"] = ""

    # =====================================
    # FEATURES
    # =====================================

    texto_final = registro["texto_reclamo"]

    registro["caps_ratio_original"] = (
        calcular_caps_ratio(texto_final)
    )

    registro["tiene_coherencia"] = (
        calcular_coherencia_basica(
            texto_final
        )
    )

    return registro


# =========================================================
# GENERAR CHUNK
# =========================================================

def generar_chunk(
    base_records,
    chunk_size,
    start_id,
):

    chunk = []

    current_id = start_id

    for _ in range(chunk_size):

        base = rc(base_records)

        registro = generar_registro(
            base,
            current_id,
        )

        chunk.append(registro)

        current_id += 1

    return chunk, current_id


# =========================================================
# MAIN
# =========================================================

def main():

    os.makedirs(LOGS_DIR, exist_ok=True)

    print("Cargando dataset base...")

    df = pd.read_csv(
        INPUT_FILE,
        sep=";",
        encoding="utf-8",
        engine="python",
    )

    df = df.loc[
        :,
        df.columns.str.strip() != ""
    ]

    df.columns = df.columns.str.strip()

    print("Columnas detectadas:")
    print(df.columns.tolist())

    # =====================================
    # TIMESTAMP
    # =====================================

    if "timestamp" not in df.columns:

        df["timestamp"] = pd.to_datetime(
            df["FECHA_PRESENTACION"],
            errors="coerce",
            dayfirst=True,
        )

    else:

        df["timestamp"] = pd.to_datetime(
            df["timestamp"],
            errors="coerce",
        )

    # =====================================
    # CAMPOS DEFAULT
    # =====================================

    if "texto_reclamo" not in df.columns:
        df["texto_reclamo"] = ""

    if "usuario_id" not in df.columns:
        df["usuario_id"] = ""

    if "ip_address" not in df.columns:
        df["ip_address"] = ""

    # =====================================
    # BASE RECORDS
    # =====================================

    base_records = df.to_dict("records")

    print(
        f"Base cargada: "
        f"{len(base_records):,}"
    )

    total_generados = 0
    next_id = 1

    first_chunk = True

    spam_count = 0
    legit_count = 0

    # =====================================
    # GENERACION
    # =====================================

    while total_generados < TARGET_SIZE:

        restante = (
            TARGET_SIZE
            - total_generados
        )

        chunk_size = min(
            CHUNK_SIZE,
            restante,
        )

        chunk, next_id = generar_chunk(
            base_records,
            chunk_size,
            next_id,
        )

        spam_chunk = sum(
            r["is_synthetic_spam"]
            for r in chunk
        )

        legit_chunk = (
            chunk_size
            - spam_chunk
        )

        spam_count += spam_chunk
        legit_count += legit_chunk

        df_chunk = pd.DataFrame(chunk)

        # =================================
        # WRITE CSV
        # =================================

        if first_chunk:

            df_chunk.to_csv(
                OUTPUT_FILE,
                index=False,
                mode="w",
            )

            first_chunk = False

        else:

            df_chunk.to_csv(
                OUTPUT_FILE,
                index=False,
                header=False,
                mode="a",
            )

        total_generados += chunk_size

        print(
            f"[OK] "
            f"{total_generados:,}/"
            f"{TARGET_SIZE:,}"
        )

        del chunk
        del df_chunk

    # =====================================
    # SUMMARY
    # =====================================

    total = total_generados

    spam_pct = (
        spam_count / total
    ) * 100

    legit_pct = (
        legit_count / total
    ) * 100

    with open(
        SUMMARY_FILE,
        "w",
        encoding="utf-8",
    ) as f:

        f.write(
            "DATASET GENERATION SUMMARY\n"
        )

        f.write(
            "==========================\n\n"
        )

        f.write(
            f"OUTPUT_FILE: "
            f"{OUTPUT_FILE}\n"
        )

        f.write(
            f"TOTAL_RECORDS: "
            f"{total:,}\n"
        )

        f.write(
            f"SPAM_RECORDS: "
            f"{spam_count:,} "
            f"({spam_pct:.2f}%)\n"
        )

        f.write(
            f"LEGITIMATE_RECORDS: "
            f"{legit_count:,} "
            f"({legit_pct:.2f}%)\n"
        )

        f.write("\nPARAMETERS\n")
        f.write("----------\n")

        f.write(
            f"TARGET_SIZE="
            f"{TARGET_SIZE}\n"
        )

        f.write(
            f"CHUNK_SIZE="
            f"{CHUNK_SIZE}\n"
        )

        f.write(
            f"SPAM_PERCENT="
            f"{SPAM_PERCENT}\n"
        )

        f.write(
            f"DIRTY_PERCENT="
            f"{DIRTY_PERCENT}\n"
        )

        f.write(
            f"NOISE_PERCENT="
            f"{NOISE_PERCENT}\n"
        )

        f.write(
            f"DUPLICATE_PERCENT="
            f"{DUPLICATE_PERCENT}\n"
        )

    print("\n===================================")
    print("DATASET GENERADO CORRECTAMENTE")
    print("===================================")

    print(f"Archivo: {OUTPUT_FILE}")
    print(f"Resumen: {SUMMARY_FILE}")


# =========================================================
# ENTRYPOINT
# =========================================================

if __name__ == "__main__":
    main()